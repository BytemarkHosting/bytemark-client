package main

import (
	"fmt"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/pathers"
	"github.com/BytemarkHosting/bytemark-client/mocks"
	"github.com/urfave/cli"
)

func TestScheduleBackups(t *testing.T) {
	type ScheduleTest struct {
		Args []string

		Name      pathers.VirtualMachineName
		DiscLabel string
		Start     string
		Interval  int

		ShouldErr  bool
		ShouldCall bool
		CreateErr  error
		BaseTestFn func(*testing.T, bool, []cli.Command) (*mocks.Config, *mocks.Client, *cli.App)
	}

	tests := []ScheduleTest{
		{
			ShouldCall: false,
			ShouldErr:  true,
			BaseTestFn: testutil.BaseTestSetup,
		},
		{
			Args:       []string{"vm-name"},
			ShouldCall: false,
			ShouldErr:  true,
			BaseTestFn: testutil.BaseTestSetup,
		},
		{
			Args:       []string{"vm-name", "disc-label"},
			Name:       pathers.VirtualMachineName{VirtualMachine: "vm-name", GroupName: pathers.GroupName{Group: "default", Account: "default-account"}},
			DiscLabel:  "disc-label",
			Start:      "00:00",
			Interval:   86400,
			ShouldCall: true,
			ShouldErr:  false,
			BaseTestFn: testutil.BaseTestAuthSetup,
		},
		{
			ShouldCall: true,
			Args:       []string{"vm-name.group.account", "disc-label", "3600"},
			Name:       pathers.VirtualMachineName{VirtualMachine: "vm-name", GroupName: pathers.GroupName{Group: "group", Account: "account"}},
			DiscLabel:  "disc-label",
			Start:      "00:00",
			Interval:   3600,
			BaseTestFn: testutil.BaseTestAuthSetup,
		},
		{
			Args:       []string{"--start", "thursday", "vm-name", "disc-label", "3235"},
			Name:       pathers.VirtualMachineName{VirtualMachine: "vm-name", GroupName: pathers.GroupName{Group: "default", Account: "default-account"}},
			DiscLabel:  "disc-label",
			Start:      "thursday",
			Interval:   3235,
			ShouldCall: true,
			ShouldErr:  true,
			CreateErr:  fmt.Errorf("intermittent failure"),
			BaseTestFn: testutil.BaseTestAuthSetup,
		},
	}

	var i int
	var test ScheduleTest
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("TestScheduleBackups #%d panicked.\r\n%v", i, r)
		}
	}()

	for i, test = range tests {
		fmt.Println(i) // fmt.Println still works even when the test panics - unlike t.Log

		config, client, app := test.BaseTestFn(t, false, commands)
		config.When("GetVirtualMachine").Return(defVM)

		retSched := brain.BackupSchedule{
			StartDate: test.Start,
			Interval:  test.Interval,
			ID:        3442,
		}
		if test.ShouldCall {
			client.When("CreateBackupSchedule", test.Name, test.DiscLabel, test.Start, test.Interval).Return(retSched, test.CreateErr).Times(1)
		} else {
			client.When("CreateBackupSchedule", test.Name, test.DiscLabel, test.Start, test.Interval).Return(retSched, test.CreateErr).Times(0)
		}
		err := app.Run(append([]string{"bytemark", "schedule", "backups"}, test.Args...))
		checkErr(t, "TestScheduleBackups", i, test.ShouldErr, err)
		verifyAndReset(t, "TestScheduleBackups", i, client)
	}
}
