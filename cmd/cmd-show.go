package cmd

import (
	"fmt"
	"strings"
)

func (cmds *CommandSet) HelpForShow() {
	// TODO(telyn): Replace instances of bigv with $0, however you get $0 in go?
	fmt.Println("bigv show")
	fmt.Println()
	fmt.Println("usage: bigv show [-j | --json] <name>")
	fmt.Println("       bigv show vm [-j | --json] <virtual machine>")
	fmt.Println("       bigv show group [-j | --json] [-v | --verbose] <group>")
	fmt.Println("       bigv show account [-j | --json] [-v | --verbose] <account>")
	fmt.Println()
	fmt.Println("Displays information about the given virtual machine, group, or account.")
	fmt.Println("If the --verbose flag is given to bigv show group or bigv show account, full details are given for each VM.")
	fmt.Println()
}

func (cmds *CommandSet) ShowVM(args []string) {
	cmds.EnsureAuth()

	name := ParseVirtualMachineName(args[0])

	vm, err := cmds.bigv.GetVirtualMachine(name)

	if err != nil {
		panic(err)
	}

	totalDiscSize := 0

	for _, disc := range vm.Discs {
		totalDiscSize += disc.Size
	}

	totalDiscSize = totalDiscSize / 1024

	// TODO(telyn): chuck this in favour of a better pretty-printer.

	fmt.Printf("= VM %s, %d cores, %d GiB RAM, %d GiB on %d discs =\r\n", vm.Name, vm.Cores, vm.Memory, totalDiscSize, len(vm.Discs))
	fmt.Printf("Hostname:		    %s\r\n", vm.Hostname)
	fmt.Printf("Hardware profile:	    %s\r\n", vm.HardwareProfile)
	fmt.Printf("Host machine:	    %s\r\n", vm.Head)
	for _, disc := range vm.Discs {
		fmt.Printf("Disc %s: %d GiB, %s grade\r\n", disc.Label, disc.Size/1024, disc.StorageGrade)
	}

	for _, nic := range vm.NetworkInterfaces {
		fmt.Printf("Network interface %s: %s\r\n", nic.Label, strings.Join(nic.Ips, ", "))

		// this is stupid
		if len(nic.ExtraIps) > 0 {

			fmt.Printf("  Extra Ips: ")
			for ip, _ := range nic.ExtraIps {
				fmt.Print(ip)
			}
			fmt.Printf("\r\n")

		}

	}

}

func (cmds *CommandSet) ShowAccount(args []string) {
	name := ParseAccountName(args[0])

	acc, err := cmds.bigv.GetAccount(name)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Account %d: %s", acc.Id, acc.Name)

}