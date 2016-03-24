package main

import (
	"bytemark.co.uk/client/cmd/bytemark/util"
	"bytemark.co.uk/client/lib"
	"bytemark.co.uk/client/util/log"
	"fmt"
	"github.com/codegangsta/cli"
	"strings"
)

func init() {
	commands = append(commands, cli.Command{
		Name: "delete",
		Description: `Deletes the given server, disc, group, account or key. Only empty groups and accounts can be deleted.
If the --purge flag is given and the target is a cloud server, will permanently delete the server. Billing will cease and you will be unable to recover the server or its data.
If the --force flag is given, you will not be prompted to confirm deletpaion.
The undelete server command may be used to restore a deleted (but not purged) server to its state prior to deletion.
`,
		Subcommands: []cli.Command{{
			Name: "account",
		}, {
			Name:    "disc",
			Aliases: []string{"disk"},
			Action: With(VirtualMachineNameProvider, DiscLabelProvider, AuthProvider, func(c *Context) (err error) {
				if !(global.Config.Force() || util.PromptYesNo("Are you sure you wish to delete this disc? It is impossible to recover.")) {
					global.Error = &util.UserRequestedExit{}
					return
				}

				err = global.Client.DeleteDisc(c.VirtualMachineName, *c.DiscLabel)
				if err != nil {
					return
				}

				return
			}),
		}, {
			Name: "group",
			Action: With(GroupProvider, func(c *Context) (err error) {
				flags := util.MakeCommonFlagSet()

				recursive := flags.Bool("recursive", false, "")

				global.Config.ImportFlags(flags)

				if len(c.Group.VirtualMachines) > 0 && *recursive {
					err = recursiveDeleteGroup(c.GroupName, c.Group)
					if err != nil {
						return
					}
				} else if !*recursive {
					err = &util.WontDeleteNonEmptyGroupError{Group: c.GroupName}
					return
				}
				err = global.Client.DeleteGroup(c.GroupName)
				return
			}),
		}, {
			Name: "key",
			Action: func(c *cli.Context) {
				user := global.Config.GetIgnoreErr("user")

				key := strings.Join(c.Args(), " ")
				if key == "" {
					log.Log("You must specify a key to delete.\r\n")
					global.Error = &util.PEBKACError{}
					return

				}

				err := EnsureAuth()
				if err != nil {
					global.Error = err
					return

				}

				err = global.Client.DeleteUserAuthorizedKey(user, key)
				if err == nil {
					log.Log("Key deleted successfullly")
					return
				} else {
					global.Error = err
					return
				}
			},
		}, {
			Name: "server",
			Action: With(VirtualMachineProvider, func(c *Context) (err error) {
				flags := util.MakeCommonFlagSet()
				purge := flags.Bool("purge", false, "")

				vm := c.VirtualMachine
				global.Config.ImportFlags(flags)

				if vm.Deleted && !*purge {
					log.Errorf("Server %s has already been deleted.\r\nIf you wish to permanently delete it, add --purge\r\n", vm.Hostname)
					// we don't return an error because we want a 0 exit code - the deletion request has happened, just not now.
					return
				}

				if !global.Config.Force() {
					fstr := fmt.Sprintf("Are you certain you wish to delete %s?", vm.Hostname)
					if *purge {
						fstr = fmt.Sprintf("Are you certain you wish to permanently delete %s? You will not be able to un-delete it.", vm.Hostname)

					}
					if !util.PromptYesNo(fstr) {
						err = &util.UserRequestedExit{}
						return

					}
				}

				err = global.Client.DeleteVirtualMachine(c.VirtualMachineName, *purge)
				if err != nil {
					return
				}
				if *purge {
					log.Logf("Server %s purged successfully.\r\n", vm.Hostname)
				} else {
					log.Logf("Server %s deleted successfully.\r\n", vm.Hostname)
				}
				return
			}),
		}},
	})
}

func recursiveDeleteGroup(name *lib.GroupName, group *lib.Group) error {
	log.Log("WARNING: The following servers will be permanently deleted, without any way to recover or un-delete them:")
	for _, vm := range group.VirtualMachines {
		log.Logf("\t%s\r\n", vm.Name)
	}
	log.Log("", "")
	if util.PromptYesNo("Are you sure you want to continue? The above servers will be permanently deleted.") {
		vmn := lib.VirtualMachineName{Group: name.Group, Account: name.Account}
		for _, vm := range group.VirtualMachines {
			vmn.VirtualMachine = vm.Name
			err := global.Client.DeleteVirtualMachine(&vmn, true)
			if err != nil {
				return err
			} else {
				log.Logf("Server %s purged successfully.\r\n", name)
			}

		}
	} else {
		return &util.UserRequestedExit{}
	}
	return nil
}

/*log.Log("usage: bytemark delete account <account>")
	log.Log("       bytemark delete disc <server> <label>")
	log.Log("       bytemark delete group [--recursive] <group>")
	//log.Log("       bytemark delete user <user>")
	log.Log("       bytemark delete key [--user=<user>] <public key identifier>")
	log.Log("       bytemark delete server [--force] [---purge] <server>")
	log.Log("       bytemark undelete server <server>")
}*/
