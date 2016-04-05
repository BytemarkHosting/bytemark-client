package main

import (
	"bufio"
	"bytemark.co.uk/client/cmd/bytemark/util"
	"bytemark.co.uk/client/lib"
	"bytemark.co.uk/client/util/log"
	"bytes"
	"encoding/json"
	"github.com/codegangsta/cli"
	"io"
	"os"
	"strings"
)

func init() {
	commands = append(commands, cli.Command{
		Name:      "debug",
		Usage:     "Test out the Bytemark API",
		UsageText: "bytemark debug [--junk-token] [--auth] [--use-billing] GET|POST|PUT|DELETE <path>",
		Description: `GET sends an HTTP GET request with an optional valid authorization header to the given path on the API endpoint and pretty-prints the received json.
The rest do similar, but PUT and POST both wait for input from stdin after authenticating. To finish entering, put an EOF (usually ctrl-d)`,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "junk-token",
				Usage: "Sets the auth token to empty - useful if you want to ensure that authenticating with credentials works, or you want to change user",
			},
			cli.BoolFlag{
				Name:  "auth",
				Usage: "Authenticate this request - without this will try to perform the call without authentication",
			},
			cli.BoolFlag{
				Name:  "use-billing",
				Usage: "Send the request to the billing endpoint instead of the brain.",
			},
		},
		Action: With(func(c *Context) error {
			shouldAuth := c.Bool("auth")

			endpoint := lib.EP_BRAIN
			if c.Bool("use-billing") {
				endpoint = lib.EP_BILLING
			}

			if c.Bool("junk-token") {
				global.Config.Set("token", "", "FLAG junk-token")
			}

			method, err := c.NextArg()
			if err != nil {
				return err
			}

			switch method {
			case "GET", "PUT", "POST", "DELETE":
				url, err := c.NextArg()

				if !strings.HasPrefix(url, "/") {
					url = "/" + url
				}
				if c.Bool("auth") {
					err := EnsureAuth()
					if err != nil {
						return err
					}
				}

				reader := io.Reader(nil)
				if method == "PUT" || method == "POST" {
					reader = bufio.NewReader(os.Stdin)
					// read until an eof
				}
				req, err := global.Client.BuildRequest(method, endpoint, url)
				if !shouldAuth {
					req, err = global.Client.BuildRequestNoAuth(method, endpoint, url)
				}
				if err != nil {
					return err
				}

				statusCode, body, err := req.Run(reader, nil)
				if err != nil {
					return err
				}
				reqURL := req.GetURL()
				log.Logf("%s %s: %d\r\n", method, reqURL.String(), statusCode)

				buf := new(bytes.Buffer)
				json.Indent(buf, body, "", "    ")
				log.Log(buf.String())
			default:
				return nil
			}
			return new(util.PEBKACError)
		}),
	})

}

// Debug makes an HTTP <method> request to the URL specified in the arguments.
// command syntax: debug <method> <url>