package main

import (
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/cheekybits/is"
	"strings"
	"testing"
)

func TestLockHWProfileCommand(t *testing.T) {
	is := is.New(t)
	config, c := baseTestAuthSetup(t, false)

	vmname := lib.VirtualMachineName{
		VirtualMachine: "test-server",
		Group:          "test-group",
		Account:        "test-account"}

	config.When("GetVirtualMachine").Return(&defVM)

	c.When("SetVirtualMachineHardwareProfileLock", &vmname, true).Return(nil).Times(1)

	err := global.App.Run(strings.Split("bytemark lock hwprofile test-server.test-group.test-account", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestUnlockHWProfileCommand(t *testing.T) {
	is := is.New(t)
	config, c := baseTestAuthSetup(t, false)

	vmname := lib.VirtualMachineName{
		VirtualMachine: "test-server",
		Group:          "test-group",
		Account:        "test-account"}

	config.When("GetVirtualMachine").Return(&defVM)

	c.When("SetVirtualMachineHardwareProfileLock", &vmname, false).Return(nil).Times(1)

	err := global.App.Run(strings.Split("bytemark unlock hwprofile test-server.test-group.test-account", " "))
	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
