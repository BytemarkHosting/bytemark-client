package main

import (
	"github.com/urfave/cli"
)

func init() {
	commands = append(commands, cli.Command{
		Name:      "hwprofiles",
		Usage:     "list hardware profiles available for cloud servers",
		UsageText: "bytemark hwprofiles [--json]",
		Description: `Hardware profiles are used by cloud servers and choosing between them can be thought of as 'which virtual motherboard should I use?'.
Generally bytemark provide two - virtio and compatibility. The virtio one has better performance but may not work with obscure operating systems, or without drivers (particularly important if you are installing windows from CD rather than our images`,
		Flags: OutputFlags("hardware profiles", "array"),
		Action: With(DefinitionsProvider, func(c *Context) error {
			return c.OutputInDesiredForm(c.Definitions.HardwareProfileDefinitions())
		}),
	}, cli.Command{
		Name:        "images",
		Aliases:     []string{"distributions", "distros"},
		Usage:       "list images available for installation on all servers",
		UsageText:   "bytemark images",
		Description: "This command lists all the images that are available for installation on Bytemark servers.",
		Flags:       OutputFlags("images", "array"),
		Action: With(DefinitionsProvider, func(c *Context) error {
			return c.OutputInDesiredForm(c.Definitions.DistributionDefinitions())
		}),
	}, cli.Command{
		Name:        "storage",
		Usage:       "list available storage grades for cloud servers",
		UsageText:   "bytemark storage",
		Description: "This outputs the available storage grades for cloud servers.",
		Flags:       OutputFlags("storage grades", "array"),
		Action: With(DefinitionsProvider, func(c *Context) error {
			return c.OutputInDesiredForm(c.Definitions.StorageGradeDefinitions())
		}),
	}, cli.Command{
		Name:        "zones",
		Usage:       "list available zones for cloud servers",
		UsageText:   "bytemark zones",
		Description: "This outputs the zones available for cloud servers to be stored and started in. Note that it is not currently possible to migrate a server between zones.",
		Flags:       OutputFlags("zones", "array"),
		Action: With(DefinitionsProvider, func(c *Context) error {
			return c.OutputInDesiredForm(c.Definitions.ZoneDefinitions())
		}),
	})
}
