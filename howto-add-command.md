How to add a command
====================

Lib
----

* Add a function to call the API to a relevant file in `lib/`. The
  files are named after what resource they work on. Basically just
  whatever the last chunk of the URL is in the API call, e.g.
  `virtual_machines`, `groups`, `accounts`. Replace `_` with `-`.

* Add the function you just made to the interface to
  `lib/interface.go`. Now `go test bigv.io/client/main` will break.

* Add a mock version of the function to `main/mocks_test.go`'s
  `mockBigVClient`. Now `go test bigv.io/client/main` will work.
  Commit if you like.

Main
----

* Add a CamelCased version of the command to `main/commands.go`'s
  `Commands` interface. Now `go test bigv.io/client/main` will break.

* Add a mock version of the function to `main/mocks_test.go`'s
  `mockCommands`. Now `go test bigv.io/client/main` will work.

* Implement the function in `main/cmd-<base>.go` where `<base>` is
  the first word of the command. See cmd-create.go and CreateVM for a
  kind of template, although CreateVM is complicated and messy.

* Add usage info to `main/cmd-<base>.go`'s `HelpFor<Base>` function.
  Split it out if you have a particularly complicated usage info. If
  you do, no need to add it to the `Commands` interface nor
  `mockCommands`.

* Add a case statement and call to the function in `Do` in
  `main/dispatcher.go`

Tests
-----

* Add a test for it to `main/dispatcher_test.go`. The test should
  ensure that your function gets called when the command is passed to
  Do.

* Add a unit test for the command to `main/commands_test.go`. You're
  trying to ensure that your function turns its arguments into the
  right parameters to the API-calling function you made at the
  beginning.

Finishing up
------------

Now do testing. `go test bigv.io/client/lib` and
`go test bigv.io/client/main`. If that works then YOU'RE NOT DONE YET

Run `make` and then try out your new `go-bigv`. Make sure the API
calls you're making actually work (--debug-level=5). Make sure errors
mostly don't cause panics.

To be honest as long as the right API calls are being made then I
(Telyn) will be able to sort out all the error handling later.

I'm intending to do get Bytemarkers to try to fuzz-test it with wacky
arguments and junky ~/.go-bigv folders and such and send me logs with
--debug-level=5 so I can squash as many as possible.

I suspect there will still be weird edge cases in the first release
and maybe for all time, but the basic idea is that `exit` in
`main/exit.go` should be smart enough 90% of the time, and all weird
errors just need wrapping in a `GoBigvWeirdError` or something.