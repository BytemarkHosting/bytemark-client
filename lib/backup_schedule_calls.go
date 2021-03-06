package lib

import (
	"strconv"

	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/pathers"
)

// CreateBackupSchedule creates a new backup schedule starting at the given date, with backups occurring every interval seconds
func (c *bytemarkClient) CreateBackupSchedule(server pathers.VirtualMachineName, discLabel string, startDate string, interval int) (sched brain.BackupSchedule, err error) {
	err = c.EnsureVirtualMachineName(&server)
	if err != nil {
		return
	}
	r, err := c.BuildRequest("POST", BrainEndpoint, "/accounts/%s/groups/%s/virtual_machines/%s/discs/%s/backup_schedules", string(server.Account), server.Group, server.VirtualMachine, discLabel)
	if err != nil {
		return
	}
	inputSchedule := brain.BackupSchedule{
		StartDate: startDate,
		Interval:  interval,
	}
	_, _, err = r.MarshalAndRun(inputSchedule, &sched)
	return
}

// DeleteBackupSchedule deletes the given backup schedule
func (c *bytemarkClient) DeleteBackupSchedule(server pathers.VirtualMachineName, discLabel string, id int) (err error) {
	err = c.EnsureVirtualMachineName(&server)
	if err != nil {
		return
	}
	r, err := c.BuildRequest("DELETE", BrainEndpoint, "/accounts/%s/groups/%s/virtual_machines/%s/discs/%s/backup_schedules/%s", string(server.Account), server.Group, server.VirtualMachine, discLabel, strconv.Itoa(id))
	if err != nil {
		return
	}
	_, _, err = r.Run(nil, nil)
	return
}
