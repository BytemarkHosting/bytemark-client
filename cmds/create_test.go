package cmds

import (
	bigv "bigv.io/client/lib"
	"bigv.io/client/mocks"
	"testing"
	//"github.com/cheekybits/is"
)

func TestCreateDiskCommand(t *testing.T) {
	c := &mocks.BigVClient{}
	config := &mocks.Config{}
	config.When("Get", "account").Return("test-account")
	config.When("Get", "token").Return("test-token")
	config.When("Force").Return(true)
	config.When("Silent").Return(true)

	config.When("ImportFlags").Return([]string{"test-vm", "archive:35"})
	name := bigv.VirtualMachineName{VirtualMachine: "test-vm"}
	c.When("ParseVirtualMachineName", "test-vm").Return(name).Times(1)
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)

	disc := bigv.Disc{Size: 35 * 1024, StorageGrade: "archive"}

	c.When("CreateDisc", name, disc).Return(nil).Times(1)

	cmds := NewCommandSet(config, c)
	cmds.CreateDiscs([]string{"test-vm", "archive:35"})

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestCreateVMCommand(t *testing.T) {
	c := &mocks.BigVClient{}
	config := &mocks.Config{}

	config.When("Get", "account").Return("test-account")
	config.When("Get", "token").Return("test-token")
	config.When("Force").Return(true)
	config.When("Silent").Return(true)
	config.When("ImportFlags").Return([]string{"test-vm"})

	c.When("ParseVirtualMachineName", "test-vm").Return(bigv.VirtualMachineName{VirtualMachine: "test-vm"})
	c.When("AuthWithToken", "test-token").Return(nil).Times(1)

	vm := bigv.VirtualMachineSpec{
		Discs: []bigv.Disc{
			bigv.Disc{
				Size:         25 * 1024,
				StorageGrade: "sata",
			},
		},
		VirtualMachine: &bigv.VirtualMachine{
			Name:                  "test-vm",
			Autoreboot:            true,
			Cores:                 1,
			Memory:                1024,
			CdromURL:              "https://example.com/example.iso",
			HardwareProfile:       "test-profile",
			HardwareProfileLocked: true,
			ZoneName:              "test-zone",
		},
		Reimage: &bigv.ImageInstall{
			Distribution: "test-image",
			RootPassword: "test-password",
		},
	}

	group := bigv.GroupName{
		Group:   "",
		Account: "",
	}

	c.When("CreateVirtualMachine", group, vm).Return(vm.VirtualMachine, nil).Times(1)

	cmds := NewCommandSet(config, c)
	cmds.CreateVM([]string{
		"--cdrom", "https://example.com/example.iso",
		"--cores", "1",
		"--discs", "25",
		"--hwprofile", "test-profile",
		"--hwprofile-locked",
		"--image", "test-image",
		"--memory", "1",
		"--root-password", "test-password",
		"--zone", "test-zone",
		"test-vm",
	})
	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}