package show

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:        "zones",
		Usage:       "show available zones for cloud servers",
		UsageText:   "show zones",
		Description: "This outputs the zones available for cloud servers to be stored and started in. Note that it is not currently possible to migrate a server between zones.",
		Flags:       app.OutputFlags("zones", "array"),
		Action: app.Action(with.Definitions, func(c *app.Context) error {
			return c.OutputInDesiredForm(c.Definitions.ZoneDefinitions(), output.List)
		}),
	})
}
