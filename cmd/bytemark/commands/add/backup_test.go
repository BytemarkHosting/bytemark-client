package add_test

import (
	"runtime/debug"
	"strings"
	"testing"
	"time"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/cheekybits/is"
	"github.com/urfave/cli"
)

func TestCreateBackup(t *testing.T) {
	is := is.New(t)
	config, c, app := testutil.BaseTestAuthSetup(t, false, commands)

	config.When("GetVirtualMachine").Return(defVM)

	vmname := lib.VirtualMachineName{
		VirtualMachine: "test-server",
		Group:          "default",
		Account:        "default-account",
	}

	c.When("CreateBackup", vmname, "test-disc").Return(brain.Backup{}, nil).Times(1)

	err := app.Run([]string{
		"bytemark", "create", "backup", "test-server", "test-disc",
	})
	is.Nil(err)
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
