package main

import (
	"reflect"
	"runtime"
	"strings"
	"testing"
	"unicode"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
)

// This test ensures that all commands have an Action, Description, Usage and UsageText
// and that all their subcommands do too.
func TestCommandsComplete(t *testing.T) {

	// TODO: Add descriptions to admin commands. it's necessary now
	t.Skip("Need to add descriptions for admin commands.")
	testutil.TraverseAllCommands(Commands(true), func(c cli.Command) {
		emptyThings := make([]string, 0, 4)
		if c.Name == "" {
			log.Log("There is a command with an empty Name.")
			t.Fail()
		}
		// if a command is only usable via its sub commands, and its usage is built from the
		// subcommands usage, its not necessary to check it.
		// incredibly hacky because this asks for the name of the method, and if it matches, just ignore it
		f := runtime.FuncForPC(reflect.ValueOf(c.Action).Pointer()).Name()
		if f == "github.com/BytemarkHosting/bytemark-client/vendor/github.com/urfave/cli.ShowSubcommandHelp" {
			return
		}

		if c.Usage == "" {
			emptyThings = append(emptyThings, "Usage")
		}
		if c.UsageText == "" {
			emptyThings = append(emptyThings, "UsageText")
		}
		if c.Description == "" {
			emptyThings = append(emptyThings, "Description")
		}
		if c.Action == nil {
			emptyThings = append(emptyThings, "Action")
		}
		if len(emptyThings) > 0 {
			t.Fail()
			log.Logf("Command %s has empty %s.\r\n", c.FullName(), strings.Join(emptyThings, ", "))
		}
	})

}

type stringPredicate func(string) bool

func checkFlagUsage(f cli.Flag, predicate stringPredicate) bool {
	switch f := f.(type) {
	case cli.BoolFlag:
		return predicate(f.Usage)
	case cli.BoolTFlag:
		return predicate(f.Usage)
	case cli.DurationFlag:
		return predicate(f.Usage)
	case cli.Float64Flag:
		return predicate(f.Usage)
	case cli.GenericFlag:
		return predicate(f.Usage)
	case cli.StringFlag:
		return predicate(f.Usage)
	case cli.StringSliceFlag:
		return predicate(f.Usage)
	}
	return false
}

func isEmpty(s string) bool {
	return s == ""
}
func notEmpty(s string) bool {
	return s != ""
}

func firstIsUpper(s string) bool {
	if s == "" {
		return false
	}

	runes := []rune(s)
	return unicode.IsUpper(runes[0])
}

func hasFullStop(s string) bool {
	return strings.Contains(s, ".")
}

func TestFlagsHaveUsage(t *testing.T) {
	testutil.TraverseAllCommands(Commands(true), func(c cli.Command) {
		for _, f := range c.Flags {
			if checkFlagUsage(f, isEmpty) {
				t.Errorf("Command %s's flag %s has empty usage\r\n", c.FullName(), f.GetName())
			}
		}
	})
	for _, flag := range app.GlobalFlags() {
		if checkFlagUsage(flag, isEmpty) {
			t.Errorf("Global flag %s doesn't have usage.", flag.GetName())
		}
	}
}

func TestUsageStyleConformance(t *testing.T) {
	testutil.TraverseAllCommandsWithContext(Commands(true), "", func(name string, c cli.Command) {
		t.Run(name, func(t *testing.T) {
			if firstIsUpper(c.Usage) {
				t.Error("Usage should be lowercase but begins with an uppercase letter")
			}

			if hasFullStop(c.Usage) {
				t.Errorf("Usage should not have full stop")
			}

			if strings.HasPrefix(c.UsageText, "bytemark ") {
				t.Error("UsageText starts with 'bytemark' - shouldn't anymore")
			}
		})
	})
}

// Tests for commands which have subcommands having the correct Description format
// the first line should start lowercase and end without a full stop, and the second
// should be blank
func TestSubcommandStyleConformance(t *testing.T) {
	testutil.TraverseAllCommands(Commands(true), func(c cli.Command) {
		if c.Subcommands == nil {
			return
		}
		if len(c.Subcommands) == 0 {
			return
		}
		if c.Description == "" {
			return
		}
		lines := strings.Split(c.Description, "\n")
		desc := []rune(lines[0])
		if unicode.IsUpper(desc[0]) {
			log.Logf("Command %s's Description begins with an uppercase letter, but it has subcommands, so should be lowercase.\r\n", c.FullName())
			t.Fail()
		}
		if strings.Contains(lines[0], ".") {
			log.Logf("The first line of Command %s's Description contains a full stop. It shouldn't.\r\n", c.FullName())
			t.Fail()
		}
		if len(lines) > 1 {
			if len(strings.TrimSpace(lines[1])) > 0 {
				log.Logf("The second line of Command %s's Description should be blank.\r\n", c.FullName())
				t.Fail()
			}
		}

	})
}
