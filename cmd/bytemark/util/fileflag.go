package util

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

//FileFlag implements the flag.Value (aka codegangsta/cli.Generic) interface
// to provide a flag value type that reads its effective value from the file named as its input.
type FileFlag struct {
	FileName string
	Value    string
}

func getPath(name string) string {
	if len(name) > 0 && name[0] == '~' {
		home := os.Getenv("HOME")
		return filepath.Join(home, name[1:])
	}
	return name
}

func (f *FileFlag) Set(name string) error {
	f.FileName = getPath(name)
	res, err := ioutil.ReadFile(f.FileName)
	f.Value = string(res)
	return err
}

func (f *FileFlag) String() string {
	return f.Value
}