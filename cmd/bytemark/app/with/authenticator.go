package with

import (
	"fmt"
	"net/url"
	"strings"

	auth3 "gitlab.bytemark.co.uk/auth/client"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/util"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/util/log"
)

type retryErr string

func (r retryErr) Error() string {
	return string(r)
}

type Authenticator struct {
	client lib.Client
	config util.ConfigManager
}

func NewAuthenticator(client lib.Client, config util.ConfigManager) Authenticator {
	return Authenticator{client: client, config: config}
}

func (a Authenticator) get2FAOTP() (otp string) {
	otp = a.config.GetIgnoreErr("2fa-otp")
	for otp == "" {
		token := util.Prompt("Enter 2FA token: ")
		a.config.Set("2fa-otp", strings.TrimSpace(token), "INTERACTION")
		otp = a.config.GetIgnoreErr("2fa-otp")
	}
	return otp
}

func (a Authenticator) tryCredentialsAttempt() error {
	credents, err := makeCredentials(a.config)

	if err != nil {
		return err
	}
	err = a.client.AuthWithCredentials(credents)

	// Handle the special case here where we just need to prompt for 2FA and try again
	if err != nil && strings.Contains(err.Error(), "Missing 2FA") {
		otp := a.get2FAOTP()

		credents["2fa"] = otp

		err = a.client.AuthWithCredentials(credents)
	}
	return err
}

func (a Authenticator) tryCredentials() (err error) {
	attempts := 3
	err = fmt.Errorf("fake error")

	for err != nil {
		attempts--
		err = a.tryCredentialsAttempt()

		if err != nil {
			if strings.Contains(err.Error(), "Badly-formed parameters") || strings.Contains(err.Error(), "Bad login credentials") {
				// if the credentials are bad in some way, make another attempt.
				if attempts <= 0 {
					return err
				}
				log.Errorf("Invalid credentials, please try again\r\n")
				a.config.Set("user", a.config.GetIgnoreErr("user"), "PRIOR INTERACTION")
				a.config.Set("pass", "", "INVALID")
				a.config.Set("yubikey-otp", "", "INVALID")
				a.config.Set("2fa-otp", "", "INVALID")
				continue
			} else {
				// if the credentials were okay and login still failed, let the user know
				return err
			}
		} else {
			// we have successfully authenticated!

			// TODO(telyn): warn on failure to write to token
			_ = a.config.SetPersistent("token", a.client.GetSessionToken(), "AUTH")

			// Check that the 2fa factor was set if --2fa-otp was specified.
			// Checking here rather than in checkFactors as it is only relevant
			// during the initial login, not subsequent validations of the
			// token (as opposed to yubikey)
			if a.config.GetIgnoreErr("2fa-otp") != "" {
				factors := a.client.GetSessionFactors()

				if a.config.GetIgnoreErr("2fa-otp") != "" {
					if !factorExists(factors, "2fa") {
						// Should never happen, as long as auth correctly returns the factors
						return fmt.Errorf("Unexpected error with 2FA login. Please report this as a bug")
					}
				}
			}
		}
	}
	return
}

func (a Authenticator) tryToken() error {
	token := a.config.GetIgnoreErr("token")

	return a.client.AuthWithToken(token)
}

func (a Authenticator) checkFactors() error {
	factors := a.client.GetSessionFactors()
	if a.config.GetIgnoreErr("yubikey") != "" {

		if a.config.GetIgnoreErr("yubikey") != "" {
			if !factorExists(factors, "yubikey") {
				// Current auth token doesn't have a yubikey,
				// so prompt the user to login again with yubikey

				// This happens when someone has logged in already,
				// but then tries to run a command with the
				// "yubikey" flag set

				a.config.Set("token", "", "FLAG yubikey")

				return EnsureAuth(a.client, a.config)
			}
		}
	}
	if a.config.GetIgnoreErr("impersonate") == "" {
		if factorExists(factors, "impersonate") {
			return retryErr("Impersonation was not requested but impersonation still happened")
		}
	} else {
		if !factorExists(factors, "impersonate") {
			return fmt.Errorf("Impersonation was requested but not achieved")
		}
	}
	return nil
}

func (a Authenticator) Authenticate() error {
	err := a.tryToken()
	if err != nil {
		// check for url.Error cause that indicates something worse than a simple auth fail.
		if aErr, ok := err.(*auth3.Error); ok {
			if _, ok := aErr.Err.(*url.Error); ok {
				return aErr
			}
		}

		log.Error("Please log in to Bytemark\r\n")

		err = a.tryCredentials()
		if err != nil {
			return err
		}
	}
	err = a.checkFactors()
	if _, ok := err.(retryErr); ok {
		fmt.Printf("%s. Retrying\n\n", err.Error())
	}
	return err
}
