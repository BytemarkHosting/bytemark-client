package admin

import (
	"fmt"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/args"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/flags"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/with"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands/admin/show"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/output"
	brainRequests "github.com/BytemarkHosting/bytemark-client/lib/requests/brain"
	"github.com/urfave/cli"
)

func init() {
	Commands = append(Commands, cli.Command{
		Name:        "show",
		Action:      cli.ShowSubcommandHelp,
		Subcommands: append(showCommands, show.Commands...),
	})
}

var showCommands = []cli.Command{
	{
		Name:      "vlans",
		Usage:     "shows available VLANs",
		UsageText: "--admin show vlans [--json]",
		Flags:     app.OutputFlags("VLANs", "array"),
		Action: app.Action(with.Auth, func(c *app.Context) error {
			vlans, err := c.Client().GetVLANs()
			if err != nil {
				return err
			}
			return c.OutputInDesiredForm(vlans, output.Table)
		}),
	},
	{
		Name:        "disc by id",
		Usage:       "displays details about a disc",
		UsageText:   "show disc by id [--json] <id>",
		Description: `Displays a collection of details about the disc.`,
		Flags: append(app.OutputFlags("disc details", "object"),
			cli.IntFlag{
				Name:  "disc",
				Usage: "the disc to display",
			},
		),
		Action: app.Action(with.Auth, args.Optional("disc"), with.RequiredFlags("disc"), func(c *app.Context) error {
			disc, err := c.Client().GetDiscByID(c.Int("disc"))
			if err != nil {
				return err
			}

			return c.OutputInDesiredForm(disc)
		}),
	},
	{
		Name:      "vlan",
		Usage:     "shows the details of a VLAN",
		UsageText: "--admin show vlan [--json] <num>",
		Flags: append(app.OutputFlags("VLAN", "object"),
			cli.IntFlag{
				Name:  "num",
				Usage: "the num of the VLAN to display",
			},
		),
		Action: app.Action(args.Optional("num"), with.RequiredFlags("num"), with.Auth, func(c *app.Context) error {
			vlan, err := c.Client().GetVLAN(c.Int("num"))
			if err != nil {
				return err
			}
			return c.OutputInDesiredForm(vlan, output.Table)
		}),
	},
	{
		Name:      "ip ranges",
		Usage:     "shows all IP ranges",
		UsageText: "--admin show ip ranges [--json]",
		Flags:     app.OutputFlags("ip ranges", "array"),
		Action: app.Action(with.Auth, func(c *app.Context) error {
			ipRanges, err := c.Client().GetIPRanges()
			if err != nil {
				return err
			}
			return c.OutputInDesiredForm(ipRanges, output.Table)
		}),
	},
	{
		Name:      "ip range",
		Usage:     "shows the details of an IP range",
		UsageText: "--admin show ip range [--json] <ip-range>",
		Flags: append(app.OutputFlags("ip range details", "object"),
			cli.StringFlag{
				Name:  "ip-range",
				Usage: "the ID or CIDR representation of the IP range to display",
			},
		),
		Action: app.Action(args.Optional("ip-range"), with.RequiredFlags("ip-range"), with.Auth, func(c *app.Context) error {
			ipRange, err := c.Client().GetIPRange(c.String("ip-range"))
			if err != nil {
				return err
			}
			return c.OutputInDesiredForm(ipRange, output.Table)
		}),
	},
	{
		Name:      "heads",
		Usage:     "shows the details of all heads",
		UsageText: "--admin show heads [--json]",
		Flags:     app.OutputFlags("heads", "array"),
		Action: app.Action(with.Auth, func(c *app.Context) error {
			heads, err := c.Client().GetHeads()
			if err != nil {
				return err
			}
			return c.OutputInDesiredForm(heads, output.Table)
		}),
	},
	{
		Name:      "head",
		Usage:     "shows the details of the specified head",
		UsageText: "--admin show head <head> [--json]",
		Flags: append(app.OutputFlags("head details", "object"),
			cli.StringFlag{
				Name:  "head",
				Usage: "the ID of the head to display",
			},
		),
		Action: app.Action(args.Optional("head"), with.RequiredFlags("head"), with.Auth, func(c *app.Context) error {
			head, err := c.Client().GetHead(c.String("head"))
			if err != nil {
				return err
			}
			return c.OutputInDesiredForm(head, output.Table)
		}),
	},
	{
		Name:      "tails",
		Usage:     "shows the details of all tails",
		UsageText: "--admin show tails [--json]",
		Flags:     app.OutputFlags("tails", "array"),
		Action: app.Action(with.Auth, func(c *app.Context) error {
			tails, err := c.Client().GetTails()
			if err != nil {
				return err
			}
			return c.OutputInDesiredForm(tails, output.Table)
		}),
	},
	{
		Name:      "tail",
		Usage:     "shows the details of the specified tail",
		UsageText: "--admin show tail <tail> [--json]",
		Flags: append(app.OutputFlags("tail details", "object"),
			cli.StringFlag{
				Name:  "tail",
				Usage: "the ID of the tail to display",
			},
		),
		Action: app.Action(args.Optional("tail"), with.RequiredFlags("tail"), with.Auth, func(c *app.Context) error {
			tail, err := c.Client().GetTail(c.String("tail"))
			if err != nil {
				return err
			}
			return c.OutputInDesiredForm(tail, output.Table)
		}),
	},
	{
		Name:      "storage pools",
		Usage:     "shows the details of all storage pools",
		UsageText: "--admin show storage pools [--json]",
		Flags:     app.OutputFlags("storage pools", "array"),
		Action: app.Action(with.Auth, func(c *app.Context) error {
			storagePools, err := c.Client().GetStoragePools()
			if err != nil {
				return err
			}
			return c.OutputInDesiredForm(storagePools, output.Table)
		}),
	},
	{
		Name:      "storage pool",
		Usage:     "shows the details of the specified storage pool",
		UsageText: "--admin show storage pool [--json] <storage-pool>",
		Flags: append(app.OutputFlags("storage pool", "object"),
			cli.StringFlag{
				Name:  "storage-pool",
				Usage: "The ID or label of the storage pool to display",
			},
		),
		Action: app.Action(args.Optional("storage-pool"), with.RequiredFlags("storage-pool"), with.Auth, func(c *app.Context) error {
			storagePool, err := c.Client().GetStoragePool(c.String("storage-pool"))
			if err != nil {
				return err
			}
			return c.OutputInDesiredForm(storagePool, output.Table)
		}),
	},
	{
		Name:      "migrating discs",
		Usage:     "shows a list of migrating discs",
		UsageText: "--admin show migrating_discs [--json]",
		Flags:     app.OutputFlags("migrating discs", "array"),
		Action: app.Action(with.Auth, func(ctx *app.Context) error {
			discs, err := ctx.Client().GetMigratingDiscs()
			if err != nil {
				return err
			}
			// this is super horrid :|
			if ctx.String("table-fields") == "" {
				err := ctx.Context.Set("table-fields", "ID, StoragePool, NewStoragePool, StorageGrade, NewStorageGrade, Size, MigrationProgress, MigrationEta, MigrationSpeed")
				if err != nil {
					return err
				}
			}
			ctx.Log("Storage sizes are in MB, speeds in MB/s, and times in seconds.")
			return ctx.OutputInDesiredForm(discs, output.Table)

		}),
	},
	{
		Name:      "migration",
		Usage:     "shows a migration job",
		UsageText: "--admin show migration [--json] <id>",
		Flags: append(app.OutputFlags("migration job", "object"),
			cli.IntFlag{
				Name:  "id",
				Usage: "the ID of the migration job",
			},
		),
		Action: app.Action(with.Auth, args.Optional("id"), with.RequiredFlags("id"), func(c *app.Context) error {
			mj, err := brainRequests.GetMigrationJob(c.Client(), c.Int("id"))
			if err != nil {
				return err
			}

			mj.Active, err = brainRequests.GetMigrationJobActiveMigrations(c.Client(), c.Int("id"))
			if err != nil {
				return err
			}

			return c.OutputInDesiredForm(mj)
		}),
	},
	{
		Name:      "migrations",
		Usage:     "shows all unfinished migration jobs",
		UsageText: "--admin show migrations",
		Action: app.Action(with.Auth, func(c *app.Context) error {
			mjs, err := brainRequests.GetMigrationJobs(c.Client())
			if err != nil {
				return err
			}

			return c.OutputInDesiredForm(mjs)
		}),
	},
	{
		Name:      "dependant servers",
		Usage:     "shows servers dependant on a head, tail or storage pool at a given time",
		UsageText: "--admin show dependant servers [--head <head>] [--tail <tail>] [--storage-pool <storage pool>] [--at <time>]",
		Flags: append(app.OutputFlags("dependant servers", "object"),
			cli.StringFlag{
				Name:  "head",
				Usage: "the ID of the head to display",
			},
			cli.StringFlag{
				Name:  "tail",
				Usage: "the ID of the tail to display",
			},
			cli.StringFlag{
				Name:  "storage-pool",
				Usage: "The ID or label of the storage pool to display",
			},
			cli.GenericFlag{
				Name:  "at",
				Value: new(flags.DateTimeFlag),
				Usage: "the date and time in history to check the dependant servers, defaults to now if unset. The flag will take most formats, e.g. '20/01/18', '15:30', '31/02/2018 18:20:33'",
			},
		),
		Action: app.Action(with.Auth, func(c *app.Context) (err error) {
			head := c.String("head")
			tail := c.String("tail")
			storage := c.String("storage-pool")
			at := c.String("at")
			var servers brain.VirtualMachines

			if head != "" {
				servers, err = brainRequests.GetServersOnHead(c.Client(), head, at)
			} else if tail != "" {
				servers, err = brainRequests.GetServersOnTail(c.Client(), tail, at)
			} else if storage != "" {
				servers, err = brainRequests.GetServersOnStoragePool(c.Client(), storage, at)
			} else {
				return fmt.Errorf("One of the following flags should be set: [--head <head>] [--tail <tail>] [--storage-pool <storage pool>]")
			}

			if err != nil {
				return
			}

			return c.OutputInDesiredForm(servers, output.Table)
		}),
	},
	{
		Name:      "dependant discs",
		Usage:     "shows discs on a head, tail or storage pool at a given time",
		UsageText: "--admin show dependant discs [--storage-pool <storage pool>] [--tail <tail>] [--at <time>]",
		Flags: append(app.OutputFlags("dependant servers", "object"),
			cli.StringFlag{
				Name:  "tail",
				Usage: "the ID of the tail to display",
			},
			cli.StringFlag{
				Name:  "storage-pool",
				Usage: "The ID or label of the storage pool to display",
			},
			cli.GenericFlag{
				Name:  "at",
				Value: new(flags.DateTimeFlag),
				Usage: "the date and time in history to check the dependant discs, defaults to now if unset. The flag will take most formats, e.g. '20/01/18', '15:30', '31/02/2018 18:20:33'",
			},
		),
		Action: app.Action(with.Auth, func(c *app.Context) (err error) {
			tail := c.String("tail")
			storage := c.String("storage-pool")
			at := c.String("at")
			var discs brain.Discs

			if tail != "" {
				discs, err = brainRequests.GetDiscsOnTail(c.Client(), tail, at)
			} else if storage != "" {
				discs, err = brainRequests.GetDiscsOnStoragePool(c.Client(), storage, at)
			} else {
				return fmt.Errorf("One of the following flags should be set: [--tail <tail>] [--storage-pool <storage pool>]")
			}

			if err != nil {
				return
			}

			return c.OutputInDesiredForm(discs, output.Table)
		}),
	},
}
