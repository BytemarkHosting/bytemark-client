package main

import (
	"fmt"
	"testing"

	"github.com/cheekybits/is"
)

func TestReapVMs(t *testing.T) {
	is := is.New(t)
	_, c, app := baseTestAuthSetup(t, true)

	c.When("ReapVMs").Return(nil).Times(1)

	err := app.Run([]string{"bytemark", "reap", "servers"})

	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestReapVMsAlias(t *testing.T) {
	is := is.New(t)
	_, c, app := baseTestAuthSetup(t, true)

	c.When("ReapVMs").Return(nil).Times(1)

	err := app.Run([]string{"bytemark", "reap", "vms"})

	is.Nil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}

func TestReapVMsError(t *testing.T) {
	is := is.New(t)
	_, c, app := baseTestAuthSetup(t, true)

	c.When("ReapVMs").Return(fmt.Errorf("Error reaping VMs")).Times(1)

	err := app.Run([]string{"bytemark", "reap", "vms"})

	is.NotNil(err)

	if ok, err := c.Verify(); !ok {
		t.Fatal(err)
	}
}
