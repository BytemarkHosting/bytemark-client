package app

import (
	"reflect"

	"github.com/urfave/cli"
)

// ProviderFunc is the function type that can be passed to Action()
type ProviderFunc func(*Context) error

// A Preprocesser is a flag.Flag that has a preprocess step that requires a Context
type Preprocesser interface {
	Preprocess(ctx *Context) error
}

// Action is a convenience function for making cli.Command.Actions that sets up a Context, runs all the providers, cleans up afterward and returns errors from the actions if there is one
func Action(providers ...ProviderFunc) func(c *cli.Context) error {
	providers = append(providers, providers[len(providers)-1])
	providers[len(providers)-2] = (*Context).Preprocess
	return func(cliContext *cli.Context) error {
		c := Context{Context: CliContextWrapper{cliContext}}
		defer cleanup(&c)

		err := foldProviders(&c, providers...)
		return err
	}
}

// Preprocess runs the Preprocess methods on all flags that implement Preprocessor
func (ctx *Context) Preprocess() error {
	if ctx.preprocessHasRun {
		return nil
	}
	ctx.Debug("Preprocessing")
	for _, flag := range ctx.Command().Flags {
		if gf, ok := flag.(cli.GenericFlag); ok {
			if pp, ok := gf.Value.(Preprocesser); ok {
				ctx.Debug("--%s b4: %#v", gf.Name, gf.Value)
				err := pp.Preprocess(ctx)
				if err != nil {
					return err
				}
				ctx.Debug("after: %#v\n", gf.Value)
			}
		}
	}
	ctx.preprocessHasRun = true
	return nil
}

// cleanup resets the value of special flags between invocations of global.App.Run so that the tests pass.
// This is needed because the init() functions are only executed once during the testing cycle.
// Outside of the tests, global.App.Run is only called once before the program closes.
func cleanup(ctx *Context) {
	allFlags := append(ctx.Command().Flags, ctx.App().Flags...)
	for _, flag := range allFlags {
		if genericFlag, ok := flag.(cli.GenericFlag); ok {
			flagValue := reflect.ValueOf(genericFlag.Value)
			if flagValue.Kind() == reflect.Ptr {
				flagValue = flagValue.Elem()
			}

			flagValue.Set(reflect.Zero(flagValue.Type()))
		}
	}
}

// foldProviders runs all the providers with the given context, stopping if there's an error
func foldProviders(ctx *Context, providers ...ProviderFunc) (err error) {
	for i, provider := range providers {
		ctx.Debug("Running provider #%d (%v)\n", i, provider)
		err = provider(ctx)
		if err != nil {
			return
		}
	}
	return
}
