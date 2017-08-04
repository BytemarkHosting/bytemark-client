package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/cheekybits/is"
)

func getFixtureDisc() brain.Disc {
	return brain.Disc{
		Label:            "",
		StorageGrade:     "sata",
		Size:             26400,
		ID:               1,
		VirtualMachineID: 1,
		StoragePool:      "fakepool",
	}
}

func getFixtureDiscSet() []brain.Disc {
	return []brain.Disc{
		getFixtureDisc(),
		brain.Disc{
			ID:           2,
			StorageGrade: "archive",
			Label:        "arch",
			Size:         1024000,
		},
		brain.Disc{
			ID:           3,
			StorageGrade: "",
			Size:         2048,
		},
	}
}

func TestLabelDisc(t *testing.T) {
	is := is.New(t)
	discs := getFixtureDiscSet()
	labelDiscs(discs)
	for _, d := range discs {
		switch d.ID {
		case 1:
			is.Equal("disc-1", d.Label)
		case 2:
			is.Equal("arch", d.Label)
		case 3:
			is.Equal("disc-3", d.Label)
		default:
			fmt.Printf("Unexpected disc ID %d\r\n", d.ID)
			t.Fail()
		}
	}
}

func TestCreateDisc(t *testing.T) {
	is := is.New(t)
	client, servers, err := mkTestClientAndServers(t, Handlers{
		brain: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path == "/accounts/account/groups/group/virtual_machines/vm/discs" && req.Method == "POST" {
				// TODO: unmarshal the disc
				// then test for sanity, equality to disk put in
				_, err := w.Write([]byte("{}"))
				if err != nil {
					t.Fatal(err)
				}
			} else if req.URL.Path == "/accounts/account/groups/group/virtual_machines/vm" {
				// TODO: return a VM that has some discs
				_, err := w.Write([]byte("{}"))
				if err != nil {
					t.Fatal(err)
				}
			} else {
				t.Fatalf("Unexpected HTTP request to %s", req.URL.String())
			}

		}),
	})
	defer servers.Close()

	is.Nil(err)
	err = client.AuthWithCredentials(map[string]string{})
	is.Nil(err)
	if err != nil {
		t.Fatal(err)
	}

	err = client.CreateDisc(VirtualMachineName{VirtualMachine: "vm", Group: "group", Account: "account"}, getFixtureDisc())

	is.Nil(err)

}

func TestDeleteDisc(t *testing.T) {
	is := is.New(t)
	client, servers, err := mkTestClientAndServers(t, Handlers{
		brain: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path == "/accounts/account/groups/group/virtual_machines/vm/discs/666" {
				if req.URL.Query().Get("purge") != "true" {
					http.NotFound(w, req)
				}

			} else if req.URL.Path == "/accounts/invalid-account" {
				http.NotFound(w, req)
			} else {
				t.Fatalf("Unexpected HTTP request to %s", req.URL.String())
			}

		}),
	})
	defer servers.Close()

	if err != nil {
		t.Fatal(err)
	}
	err = client.AuthWithCredentials(map[string]string{})
	if err != nil {
		t.Fatal(err)
	}

	err = client.DeleteDisc(VirtualMachineName{VirtualMachine: "vm", Group: "group", Account: "account"}, "666")
	is.Nil(err)

}

func TestResizeDisc(t *testing.T) {
	is := is.New(t)
	expectedDisc := map[string]interface{}{
		"size": 35.0, // json library treats all numbers as float64
	}
	client, servers, err := mkTestClientAndServers(t, Handlers{
		brain: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path == "/accounts/account/groups/group/virtual_machines/vm/discs/666" {
				bytes, err := ioutil.ReadAll(req.Body)
				is.Nil(err)
				var actualDisc map[string]interface{}
				err = json.Unmarshal(bytes, &actualDisc)
				is.Nil(err)
				// TODO check that we send the right stuff
				if !reflect.DeepEqual(actualDisc, expectedDisc) {
					t.Errorf("Resize disc request wasn't as expected.\r\nExpected: %#v\r\nActual: %#v", expectedDisc, actualDisc)
				}
			} else if req.URL.Path == "/accounts/invalid-account" {
				http.NotFound(w, req)
			} else {
				t.Fatalf("Unexpected HTTP request to %s", req.URL.String())
			}

		}),
	})
	defer servers.Close()

	if err != nil {
		t.Fatal(err)
	}
	err = client.AuthWithCredentials(map[string]string{})
	if err != nil {
		t.Fatal(err)
	}

	err = client.ResizeDisc(VirtualMachineName{VirtualMachine: "vm", Group: "group", Account: "account"}, "666", 35)
	is.Nil(err)

}

func TestShowDisc(t *testing.T) {
	is := is.New(t)

	client, servers, err := mkTestClientAndServers(t, Handlers{
		brain: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path == "/accounts/account/groups/group/virtual_machines/vm/discs/666" {
				bytes, err := json.Marshal(getFixtureDisc())
				is.Nil(err)
				_, err = w.Write(bytes)
				if err != nil {
					t.Fatal(err)
				}
			} else if req.URL.Path == "/accounts/invalid-account" {
				http.NotFound(w, req)
			} else {
				t.Fatalf("Unexpected HTTP request to %s", req.URL.String())
			}

		}),
	})
	defer servers.Close()

	if err != nil {
		t.Fatal(err)
	}
	err = client.AuthWithCredentials(map[string]string{})
	if err != nil {
		t.Fatal(err)
	}

	disc, err := client.GetDisc(VirtualMachineName{VirtualMachine: "vm", Group: "group", Account: "account"}, "666")
	if err != nil {
		t.Fatal(err)
	}
	is.Nil(err)
	fx := getFixtureDisc()

	is.Equal(fx.ID, disc.ID)
	is.Equal(fx.Label, disc.Label)
	is.Equal(fx.Size, disc.Size)
	is.Equal(fx.StorageGrade, disc.StorageGrade)
	is.Equal(fx.StoragePool, disc.StoragePool)

}
