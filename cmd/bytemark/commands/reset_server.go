package commands

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/flags"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:        "reset",
		Usage:       "restart a server as though the reset button had been pushed",
		UsageText:   "reset server <server>",
		Description: "For cloud servers, this does not cause the qemu process to be restarted. This means that the server will remain on the same head and will not notice hardware changes.",
		Action:      cli.ShowSubcommandHelp,
		Flags: []cli.Flag{
			cli.GenericFlag{
				Name:  "server",
				Usage: "the server to reset",
				Value: new(flags.VirtualMachineNameFlag),
			},
		},
		Subcommands: []cli.Command{{
			Name:        "server",
			Usage:       "restart a server as though the reset button had been pushed",
			UsageText:   "reset server <server>",
			Description: "For cloud servers, this does not cause the qemu process to be restarted. This means that the server will remain on the same head and will not notice hardware changes.",
			Flags: []cli.Flag{
				cli.GenericFlag{
					Name:  "server",
					Usage: "the server to reset",
					Value: new(flags.VirtualMachineNameFlag),
				},
			},
			Action: app.Action(args.Optional("server"), with.RequiredFlags("server"), with.Auth, func(c *app.Context) (err error) {
				vmName := flags.VirtualMachineName(c, "server")
				c.LogErr("Attempting to reset %v...\r\n", vmName)
				err = c.Client().ResetVirtualMachine(vmName)
				if err != nil {
					return err
				}

				c.LogErr("%v reset successfully.\r\n", vmName)
				return
			}),
		}},
	})
}
