package config

import (
	"flag"

	"github.com/BytemarkHosting/bytemark-client/lib/pathers"
)

// Manager is an interface defining a key->value store that also knows where the values were set from.
type Manager interface {
	Get(string) (string, error)
	GetIgnoreErr(string) string
	GetBool(string) (bool, error)
	GetV(string) (Var, error)
	GetSessionValidity() (int, error)
	GetVirtualMachine() pathers.VirtualMachineName
	GetGroup() pathers.GroupName
	GetAll() (Vars, error)
	Set(string, string, string)
	SetPersistent(varname string, value string, source string) error
	Unset(string) error
	GetDebugLevel() int
	EndpointName() string
	PanelURL() string
	ConfigDir() string

	ImportFlags(*flag.FlagSet) []string
}
