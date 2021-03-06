package app

import (
	"io"
	"os"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/cliutil"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/config"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/util/log"
	"github.com/urfave/cli"
)

// BaseAppSetup sets up a cli.App for the given commands and config
func BaseAppSetup(flags []cli.Flag, commands []cli.Command) (app *cli.App, err error) {
	app = cli.NewApp()
	app.Version = lib.Version
	app.Flags = flags
	app.Commands = commands
	app.Usage = "Command-line interface to Bytemark Cloud services"
	app.Writer = io.MultiWriter(
		log.LogFile,
		os.Stdout,
	)
	app.ErrWriter = io.MultiWriter(
		log.LogFile,
		os.Stderr,
	)

	app.Commands = cliutil.CreateMultiwordCommands(app.Commands)
	return

}

// SetClientAndConfig adds the client and config to the given app.
// it abstracts away setting the Metadata on the app. Mostly so that we get some type-checking.
// without it - it's just assigning to an interface{} which will always succeed,
// and which would near-inevitably result in hard-to-debug null pointer errors down the line.
func SetClientAndConfig(app *cli.App, client lib.Client, config config.Manager) {
	if app.Metadata == nil {
		app.Metadata = make(map[string]interface{})
	}
	app.Metadata["client"] = client
	app.Metadata["config"] = config
	app.Metadata["prompter"] = util.NewPrompter()
}

// SetPrompter sets the prompter for the given app. It does not normally need
// to be set (as this is done for you when calling SetClientAndConfig in
// cmd/bytemark.main() )
// This is used by command tests such as TestDeleteGroup in
// cmd/bytemark/commands/delete/group_test.go to allow for mocking user input.
// See that test, or cmd/bytemark/app/auth/authenticator_test.go for example
// usage of a mock Prompter.
func SetPrompter(app *cli.App, prompter util.Prompter) {
	if app.Metadata == nil {
		app.Metadata = make(map[string]interface{})
	}
	app.Metadata["prompter"] = prompter
}
