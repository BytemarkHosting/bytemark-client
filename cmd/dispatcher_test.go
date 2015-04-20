package cmd

import (
	"testing"
	//"github.com/cheekybits/is"
)

func doDispatchTest(t *testing.T, config *mockConfig, commands *mockCommands, args ...string) {
	d := NewDispatcherWithCommands(config, commands)

	if args == nil {
		args = []string{}
	}

	d.Do(args)

	if ok, err := commands.Verify(); !ok {
		t.Fatalf("Test with args %v failed: %v", args, err)
	}
}

func TestDispatchDoDebug(t *testing.T) {
	commands := &mockCommands{}
	config := &mockConfig{}
	config.When("Get", "endpoint").Return("endpoint.example.com")
	config.When("GetDebugLevel").Return(0)

	commands.When("Debug", []string{"GET", "/test"}).Times(1)
	doDispatchTest(t, config, commands, "debug", "GET", "/test")

	commands.Reset()

}

func TestDispatchDoHelp(t *testing.T) {
	commands := &mockCommands{}
	config := &mockConfig{}
	config.When("Get", "endpoint").Return("endpoint.example.com")
	config.When("GetDebugLevel").Return(0)

	commands.When("Help", []string{}).Times(1)

	doDispatchTest(t, config, commands)

	commands.Reset()
	commands.When("Help", []string{}).Times(1)
	doDispatchTest(t, config, commands, "help")

	commands.When("Help", []string{"show"}).Times(1)
	doDispatchTest(t, config, commands, "help", "show")
}

func TestDispatchDoConfig(t *testing.T) {
	commands := &mockCommands{}
	config := &mockConfig{}
	config.When("Get", "endpoint").Return("endpoint.example.com")
	config.When("GetDebugLevel").Return(0)

	commands.When("Config", []string{}).Times(1)
	doDispatchTest(t, config, commands, "config")

	commands.Reset()
	commands.When("Config", []string{"set"}).Times(1)
	doDispatchTest(t, config, commands, "config", "set")

	commands.Reset()
	commands.When("Config", []string{"set", "variablename"}).Times(1)
	doDispatchTest(t, config, commands, "config", "set", "variablename")

	commands.Reset()
	commands.When("Config", []string{"set", "variablename", "value"}).Times(1)
	doDispatchTest(t, config, commands, "config", "set", "variablename", "value")
}

//
//func TestDispatchDoShow(t *testing.T) {
//	commands := &mockCommands{}
//	config := &mockConfig{}
//	config.When("Get", "endpoint").Return("endpoint.example.com")
//	config.When("GetDebugLevel").Return(0)
//
//	doDispatchTest(t, config, commands)
//}
//
//func TestDispatchDoUnset(t *testing.T) {
//	commands := &mockCommands{}
//	config := &mockConfig{}
//	config.When("Get", "endpoint").Return("endpoint.example.com")
//	config.When("GetDebugLevel").Return(0)
//
//	doDispatchTest(t, config, commands)
//}
