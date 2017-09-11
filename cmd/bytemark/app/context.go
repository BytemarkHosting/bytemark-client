package app

import (
	"fmt"
	"io"
	"net"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/urfave/cli"
)

// Context is a wrapper around urfave/cli.Context which provides easy access to
// the next unused argument and can have various bytemarky types attached to it
// in order to keep code DRY
type Context struct {
	Context        innerContext
	Account        *lib.Account
	Authed         bool
	Definitions    *lib.Definitions
	Disc           *brain.Disc
	Group          *brain.Group
	Privilege      brain.Privilege
	User           *brain.User
	VirtualMachine *brain.VirtualMachine

	currentArgIndex int
	preproDone      bool
}

// Reset replaces the Context with a blank one (keeping the cli.Context)
func (c *Context) Reset() {
	*c = Context{
		Context: c.Context,
	}
}

// App returns the cli.App that this context is part of. Usually this will be the same as global.App, but it's nice to depend less on globals.
func (c *Context) App() *cli.App {
	return c.Context.App()
}

// args returns all the args that were passed to the Context (i.e. all the args passed to this (sub)command)
func (c *Context) args() cli.Args {
	return c.Context.Args()
}

// Args returns all the unused arguments
func (c *Context) Args() []string {
	return c.args()[c.currentArgIndex:]
}

// Debug runs fmt.Fprintf on the args, outputting to the App's debugWriter.
// In tests, this is a TestWriter. Otherwise it's nil for now - but might be
// changed to the debug.log File in the future.
func (c *Context) Debug(format string, values ...interface{}) {
	dw, ok := c.App().Metadata["debugWriter"]
	if !ok {
		return
	}
	if wr, ok := dw.(io.Writer); ok {
		fmt.Fprintf(wr, format, values...)
	}
}

// Command returns the cli.Command this context is for
func (c *Context) Command() cli.Command {
	return c.Context.Command()
}

// Config returns the config attached to the App this Context is for
func (c *Context) Config() util.ConfigManager {
	if config, ok := c.App().Metadata["config"].(util.ConfigManager); ok {
		return config
	}
	return nil
}

// Client returns the API client attached to the App this Context is for
func (c *Context) Client() lib.Client {
	if client, ok := c.App().Metadata["client"].(lib.Client); ok {
		return client
	}
	return nil
}

// NextArg returns the next unused argument, and marks it as used.
func (c *Context) NextArg() (string, error) {
	if len(c.args()) <= c.currentArgIndex {
		return "", c.Help("not enough arguments were specified")
	}
	arg := c.args()[c.currentArgIndex]
	c.currentArgIndex++
	return arg, nil
}

// Help creates a UsageDisplayedError that will output the issue and a message to consult the documentation
func (c *Context) Help(whatsyourproblem string) (err error) {
	return util.UsageDisplayedError{TheProblem: whatsyourproblem, Command: c.Command().FullName()}
}

// flags below

// Bool returns the value of the named flag as a bool
func (c *Context) Bool(flagname string) bool {
	return c.Context.Bool(flagname)
}

// Discs returns the discs passed along as the named flag.
// I can't imagine why I'd ever name a disc flag anything other than --disc, but the flexibility is there just in case.
func (c *Context) Discs(flagname string) []brain.Disc {
	disc, ok := c.Context.Generic(flagname).(*util.DiscSpecFlag)
	if ok {
		return []brain.Disc(*disc)
	}
	return []brain.Disc{}
}

// FileName returns the name of the file given by the named flag
func (c *Context) FileName(flagname string) string {
	file, ok := c.Context.Generic(flagname).(*util.FileFlag)
	if ok {
		return file.FileName
	}
	return ""
}

// FileContents returns the contents of the file given by the named flag.
func (c *Context) FileContents(flagname string) string {
	file, ok := c.Context.Generic(flagname).(*util.FileFlag)
	if ok {
		return file.Value
	}
	return ""
}

// GroupName returns the named flag as a lib.GroupName
func (c *Context) GroupName(flagname string) (gp lib.GroupName) {
	gpNameFlag, ok := c.Context.Generic(flagname).(*GroupNameFlag)
	if !ok {
		fmt.Println("WRONGO")
		return lib.GroupName{}
	}
	if gpNameFlag.GroupName == nil {
		fmt.Println("NILO")
		return lib.GroupName{}
	}
	return *gpNameFlag.GroupName
}

// Int returns the value of the named flag as an int
func (c *Context) Int(flagname string) int {
	return c.Context.Int(flagname)
}

// Int64 returns the value of the named flag as an int64
func (c *Context) Int64(flagname string) int64 {
	return c.Context.Int64(flagname)
}

// IPs returns the ips passed along as the named flag.
func (c *Context) IPs(flagname string) []net.IP {
	ips, ok := c.Context.Generic(flagname).(*util.IPFlag)
	if ok {
		return []net.IP(*ips)
	}
	return []net.IP{}
}

// PrivilegeFlag returns the named flag as a PrivilegeFlag
func (c *Context) PrivilegeFlag(flagname string) PrivilegeFlag {
	priv, ok := c.Context.Generic(flagname).(*PrivilegeFlag)
	if ok {
		return *priv
	}
	return PrivilegeFlag{}
}

// String returns the value of the named flag as a string
func (c *Context) String(flagname string) string {
	if c.Context.IsSet(flagname) || c.Context.String(flagname) != "" {
		c.Debug("IsSet || String() != nil")
		return c.Context.String(flagname)
	}
	return c.Context.GlobalString(flagname)
}

// Size returns the value of the named SizeSpecFlag as an int in megabytes
func (c *Context) Size(flagname string) int {
	size, ok := c.Context.Generic(flagname).(*util.SizeSpecFlag)
	if ok {
		return int(*size)
	}
	return 0
}

// ResizeFlag returns the named ResizeFlag
func (c *Context) ResizeFlag(flagname string) ResizeFlag {
	size, ok := c.Context.Generic(flagname).(*ResizeFlag)
	if ok {
		return *size
	}
	return ResizeFlag{}
}

// VirtualMachineName returns the named flag as a lib.VirtualMachineName
func (c *Context) VirtualMachineName(flagname string) (vm lib.VirtualMachineName) {
	vmNameFlag, ok := c.Context.Generic(flagname).(*VirtualMachineNameFlag)
	if !ok {
		fmt.Println("WRONGO")
		return c.Config().GetVirtualMachine()
	}
	if vmNameFlag.VirtualMachineName == nil {
		fmt.Println("NILO")
		return lib.VirtualMachineName{}
	}

	return *vmNameFlag.VirtualMachineName
}