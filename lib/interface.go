package lib

import (
	auth3 "bytemark.co.uk/auth3/client"
	"net/http"
)

type Client interface {
	// Getters
	GetEndpoint() string
	GetSessionToken() string

	// Setters
	SetDebugLevel(int)

	AuthWithToken(string) error
	AuthWithCredentials(auth3.Credentials) error

	RequestAndUnmarshal(auth bool, method, path, requestBody string, output interface{}) error
	RequestAndRead(auth bool, method, path, requestBody string) (responseBody []byte, err error)
	Request(auth bool, method string, location string, requestBody string) (req *http.Request, res *http.Response, err error)

	// ACCOUNTS
	GetAccount(name string) (account *Account, err error)

	// GROUPS

	// VIRTUAL MACHINES
	GetVirtualMachine(name VirtualMachineName) (vm *VirtualMachine, err error)
}