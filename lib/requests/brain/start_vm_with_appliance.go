package brain

import (
	"bytes"
	"fmt"

	"github.com/BytemarkHosting/bytemark-client/lib"
)

// StartVirtualMachineWithAppliance starts the named virtual machine using the named appliance.
// returns nil on success or an error otherwise.
func StartVirtualMachineWithAppliance(client lib.Client, vmName lib.VirtualMachineName, applianceName string) (err error) {
	err = client.EnsureVirtualMachineName(&vmName)
	if err != nil {
		return err
	}
	r, err := client.BuildRequest("PUT", lib.BrainEndpoint, "/accounts/%s/groups/%s/virtual_machines/%s", vmName.Account, vmName.Group, vmName.VirtualMachine)
	if err != nil {
		return err
	}

	update := map[string]interface{}{
		"autoreboot_on": true,
		"power_on":      true,
		"appliance": map[string]interface{}{
			"name":      applianceName,
			"permanent": false,
		},
	}
	_, _, err = r.MarshalAndRun(update, nil)

	return err
}
