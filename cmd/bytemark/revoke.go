package main

import (
	"fmt"

	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
)

func init() {
	commands = append(commands, cli.Command{
		Name:        "revoke",
		Usage:       "revoke privileges on bytemark self-service objects from other users",
		UsageText:   "bytemark revoke <privilege> [on] <object> [from] <user>\r\nbytemark grant cluster_admin [to] <user>",
		Description: "Revoke a privilege from a user for a particular bytemark object\r\n\r\n" + privilegeText,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "yubikey-required",
				Usage: "Set if the privilege should require a yubikey.",
			},
			cli.GenericFlag{
				Name:  "privilege",
				Usage: "the privilege to revoke",
				Value: new(PrivilegeFlag),
			},
		},
		Action: With(JoinArgs("privilege"), RequiredFlags("privilege"), PrivilegeProvider("privilege"), func(c *Context) (err error) {
			pf := c.PrivilegeFlag("privilege")
			c.Privilege.YubikeyRequired = c.Bool("yubikey-required")

			var privs brain.Privileges
			switch c.Privilege.TargetType() {
			case brain.PrivilegeTargetTypeVM:
				privs, err = global.Client.GetPrivilegesForVirtualMachine(*pf.VirtualMachineName)
				if err != nil {
					return
				}
			case brain.PrivilegeTargetTypeGroup:
				privs, err = global.Client.GetPrivilegesForGroup(*pf.GroupName)
				if err != nil {
					return
				}
			case brain.PrivilegeTargetTypeAccount:
				privs, err = global.Client.GetPrivilegesForAccount(pf.AccountName)
				if err != nil {
					return
				}
			default:
				privs, err = global.Client.GetPrivileges(pf.Username)
				if err != nil {
					return
				}
			}
			i := privs.IndexOf(c.Privilege)
			if i == -1 {
				return fmt.Errorf("Couldn't find such a privilege to revoke")
			}

			err = global.Client.RevokePrivilege(privs[i])
			if err == nil {
				log.Outputf("Revoked %s\r\n", c.PrivilegeFlag("privilege").String())

			}
			return
		}),
	})
}
