package commands_test

import (
	"strings"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib/pathers"
	"github.com/cheekybits/is"
)

func TestResetCommand(t *testing.T) {
	is := is.New(t)
	config, c, app := testutil.BaseTestAuthSetup(t, false, commands.Commands)
	vmn := pathers.VirtualMachineName{VirtualMachine: "test-server", GroupName: pathers.GroupName{Group: "test-group", Account: "test-account"}}

	config.When("GetVirtualMachine").Return(pathers.VirtualMachineName{VirtualMachine: "smog-daniels", GroupName: pathers.GroupName{Group: "creation", Account: "science"}})

	c.When("ResetVirtualMachine", vmn).Times(1)

	err := app.Run(strings.Split("bytemark reset server test-server.test-group.test-account", " "))
	is.Nil(err)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
