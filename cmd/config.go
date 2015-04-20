package cmd

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var CONFIG_VARIABLES = [...]string{
	"endpoint",
	"auth-endpoint",
	"user",
	"account",
	"token",
	"debug-level",
}

// ConfigVar is a struct which contains a name-value-source triplet
// Source is up to two words separated by a space. The first word is the source type: FLAG, ENV, DIR, CODE.
// The second is the name of the flag/file/environment var used.
type ConfigVar struct {
	Name   string
	Value  string
	Source string
}

type ConfigManager interface {
	Get(string) string
	GetV(string) ConfigVar
	GetAll() []ConfigVar
	Set(string, string, string)
	SetPersistent(string, string, string)
	Unset(string)
	GetDebugLevel() int
}

// Params currently used:
// token - an OAuth 2.0 bearer token to use when authenticating
// username - the default username to use - if not present, $USER
// endpoint - the default endpoint to use - if not present, https://uk0.bigv.io
// auth-endpoint - the default auth API endpoint to use - if not present, https://auth.bytemark.co.uk

// A Config determines the configuration of the bigv client.
// It's responsible for handling things like the credentials to use and what endpoints to talk to.
//
// Each configuration item is read from the following places, falling back to successive places:
//
// Per-command command-line flags, global command-line flags, environment variables, configuration directory, hard-coded defaults
//
//The location of the configuration directory is read from global command-line flags, or is otherwise ~/.go-bigv
//
type Config struct {
	debugLevel  int
	Dir         string
	Memo        map[string]ConfigVar
	Definitions map[string]string
}

// Do I really need to have the flags passed in here?
// Yes. Doing commands will be sorted out in a different place, and I don't want to touch it here.

// NewConfig sets up a new config struct. Pass in an empty string to default to ~/.go-bigv
func NewConfig(configDir string, flags *flag.FlagSet) (config *Config) {
	config = new(Config)
	config.Memo = make(map[string]ConfigVar)
	config.Dir = filepath.Join(os.Getenv("HOME"), "/.go-bigv")
	if os.Getenv("BIGV_CONFIG_DIR") != "" {
		config.Dir = os.Getenv("BIGV_CONFIG_DIR")
	}

	if configDir != "" {
		config.Dir = configDir
	}

	err := os.MkdirAll(config.Dir, 0700)
	if err != nil {
		// TODO(telyn): Better error handling here

		panic(err)
	}

	stat, err := os.Stat(config.Dir)
	if err != nil {
		// TODO(telyn): Better error handling here
		panic(err)
	}

	if !stat.IsDir() {
		fmt.Printf("%s is not a directory", config.Dir)
		panic("Cannot continue")
	}

	if flags != nil {
		// dump all the flags into the memo
		// should be reet...reet?
		flags.Visit(func(f *flag.Flag) {
			config.Memo[f.Name] = ConfigVar{
				f.Name,
				f.Value.String(),
				"FLAG " + f.Name,
			}
		})
	}
	debugLevel, err := strconv.ParseInt(config.Get("debug-level"), 10, 0)
	if err == nil {
		config.debugLevel = int(debugLevel)
	}
	return config
}

func (config *Config) GetDebugLevel() int {
	return config.debugLevel
}

// GetPath joins the given string onto the end of the Config.Dir path
func (config *Config) GetPath(name string) string {
	return filepath.Join(config.Dir, name)
}
func (config *Config) GetUrl(path ...string) *url.URL {
	url, err := url.Parse(config.Get("endpoint"))
	if err != nil {
		panic("Endpoint is not a valid URL")
	}
	url.Parse("/" + strings.Join([]string(path), "/"))
	return url
}

func (config *Config) LoadDefinitions() {
	stat, err := os.Stat(config.GetPath("definitions"))

	if err != nil || time.Since(stat.ModTime()) > 24*time.Hour {
		// TODO(telyn): grab it off the internet
		//		url := config.GetUrl("definitions.json")
	} else {
		_, err := ioutil.ReadFile(config.GetPath("definitions"))
		if err != nil {
			panic("Couldn't load definitions")
		}
	}

}

func (config *Config) Get(name string) string {
	return config.GetV(name).Value
}

func (config *Config) GetV(name string) ConfigVar {
	// try to read the Memo
	name = strings.ToLower(name)
	if val, ok := config.Memo[name]; ok {
		return val
	} else {
		return config.read(name)
	}
}

func (config *Config) GetAll() []ConfigVar {
	vars := make([]ConfigVar, 4)
	for _, v := range CONFIG_VARIABLES {
		vars = append(vars, config.GetV(v))
	}
	return vars
}

func (config *Config) GetDefault(name string) ConfigVar {
	// ideally most of these should just be	os.Getenv("BIGV_"+name.Upcase().Replace("-","_"))
	switch name {
	case "user":
		// we don't actually want to default to USER - that will happen during Dispatcher's PromptForCredentials so it can be all "Hey you should bigv config set user <youruser>"
		return ConfigVar{"user", os.Getenv("BIGV_USER"), "ENV BIGV_USER"}
	case "endpoint":
		v := ConfigVar{"endpoint", "https://uk0.bigv.io", "CODE"}

		val := os.Getenv("BIGV_ENDPOINT")
		if val != "" {
			v.Value = val
			v.Source = "ENV BIGV_ENDPOINT"
		}
		return v
	case "auth-endpoint":
		v := ConfigVar{"endpoint", "https://auth.bytemark.co.uk", "CODE"}

		val := os.Getenv("BIGV_AUTH_ENDPOINT")
		if val != "" {
			v.Value = val
			v.Source = "ENV BIGV_AUTH_ENDPOINT"
		}
		return v
	case "account":
		val := os.Getenv("BIGV_ACCOUNT")
		if val != "" {
			return ConfigVar{
				"account",
				val,
				"ENV BIGV_AUTH_ENDPOINT",
			}
		}
		return config.GetDefault("user")
	case "debug-level":
		v := ConfigVar{"debug-level", "0", "CODE"}
		if val := os.Getenv("BIGV_DEBUG_LEVEL"); val != "" {
			v.Value = val
		}
		return v
	}
	return ConfigVar{"", "", ""}
}

func (config *Config) read(name string) ConfigVar {
	contents, err := ioutil.ReadFile(config.GetPath(name))
	if err != nil {
		if os.IsNotExist(err) {
			return config.GetDefault(name)
		}
		fmt.Printf("Couldn't read config for %s", name)
		panic(err)
	}

	return ConfigVar{name, string(contents), "FILE " + config.GetPath(name)}
}

// Set stores the given key-value pair in config's Memo. This storage does not persist once the program terminates.
func (config *Config) Set(name, value, source string) {
	config.Memo[name] = ConfigVar{name, value, source}
}

// SetPersistent writes a file to the config directory for the given key-value pair.
func (config *Config) SetPersistent(name, value, source string) {
	config.Set(name, value, source)
	err := ioutil.WriteFile(config.GetPath(name), []byte(value), 0600)
	if err != nil {
		fmt.Println("Couldn't write to config directory " + config.Dir)
		panic(err)
	}
}

// Unset removes the named key from both config's Memo and the user's config directory.
func (config *Config) Unset(name string) {
	delete(config.Memo, name)
	// TODO(telyn): handle errors here or don't
	os.Remove(config.GetPath(name))
}
