package migrate_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands/admin"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/pathers"
	"github.com/BytemarkHosting/bytemark-client/mocks"
	"github.com/urfave/cli"
)

func TestMigrateServer(t *testing.T) {
	tests := []struct {
		name       string
		args       string
		head       string
		vm         pathers.VirtualMachineName
		migrateErr error
		shouldErr  bool
	}{{
		name: "with new head",
		head: "stg-h1",
		vm: pathers.VirtualMachineName{
			VirtualMachine: "vm123",
			GroupName: pathers.GroupName{
				Group:   "group",
				Account: "account",
			},
		},

		args: "migrate vm vm123.group.account stg-h1",
	}, {
		name: "without new head",
		vm: pathers.VirtualMachineName{
			VirtualMachine: "vm122",
			GroupName: pathers.GroupName{
				Group:   "group",
				Account: "account",
			},
		},
		args: "migrate vm vm122.group.account",
	}, {
		name: "error",
		vm: pathers.VirtualMachineName{
			VirtualMachine: "vm121",
			GroupName: pathers.GroupName{
				Group:   "group",
				Account: "account",
			},
		},
		args:       "migrate vm vm121.group.account",
		migrateErr: fmt.Errorf("all the heads are down oh no"),
		shouldErr:  true,
	}, {
		name: "id",
		vm: pathers.VirtualMachineName{
			VirtualMachine: "1123",
			GroupName: pathers.GroupName{
				Group:   "test-group",
				Account: "test-account",
			},
		},
		args: "migrate vm 1123",
	}}

	for _, test := range tests {
		ct := testutil.CommandT{
			Name:      test.name,
			Auth:      true,
			Admin:     true,
			Args:      test.args,
			ShouldErr: test.shouldErr,
			Commands:  admin.Commands,
		}
		if !test.shouldErr {
			ct.OutputMustMatch = []*regexp.Regexp{
				// ensure that the output contains the truncated hostname
				// of the vm returned from GetVirtualMachineName below - in
				// other words, the 'real' name of the VM that the user wanted
				// to migrate.
				regexp.MustCompile(`real\.cool\.vm`),
			}
		}
		ct.Run(t, func(t *testing.T, config *mocks.Config, client *mocks.Client, app *cli.App) {
			config.When("GetVirtualMachine").Return(pathers.VirtualMachineName{
				VirtualMachine: "test-vm",
				GroupName: pathers.GroupName{
					Group:   "test-group",
					Account: "test-account",
				},
			})
			client.When("GetVirtualMachine", test.vm).Return(brain.VirtualMachine{
				ID:       1123,
				Hostname: "real.cool.vm.uk0.bigv.io",
			})

			client.When("MigrateVirtualMachine", test.vm, test.head).Return(test.migrateErr).Times(1)
		})
	}
}
