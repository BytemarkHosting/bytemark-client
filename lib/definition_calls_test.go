package lib_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil"
	"github.com/cheekybits/is"
)

func TestProcessDefinitions(t *testing.T) {
	is := is.New(t)
	var defsJS lib.JSONDefinitions
	err := json.Unmarshal([]byte(fixtureDefinitionsJSON), &defsJS)

	is.Nil(err)
	if err != nil {
		fmt.Print(err)
		return
	}

	defs := defsJS.Process()

	is.Equal(3, len(defs.StorageGrades))
	is.Equal(2, len(defs.HardwareProfiles))
	is.Equal(12, len(defs.Distributions))

	for k, l := range fixtureDefinitions.Distributions {
		is.Equal(l, defs.Distributions[k])
	}
	for k, l := range fixtureDefinitions.StorageGrades {
		is.Equal(l, defs.StorageGrades[k])
	}
	for k, l := range fixtureDefinitions.HardwareProfiles {
		is.Equal(l, defs.HardwareProfiles[k])
	}
	for k, l := range fixtureDefinitions.ZoneNames {
		is.Equal(l, defs.ZoneNames[k])
	}

}

func TestReadDefinitions(t *testing.T) {
	is := is.New(t)

	rts := testutil.RequestTestSpec{
		Method:   "GET",
		Endpoint: lib.BrainEndpoint,
		URL:      "/definitions",
		Response: json.RawMessage(fixtureDefinitionsJSON),
	}

	rts.Run(t, testutil.Name(0), false, func(client lib.Client) {
		defs, err := client.ReadDefinitions()
		if err != nil {
			t.Fatal(err)
		}

		is.Equal(3, len(defs.StorageGrades))
		is.Equal(2, len(defs.HardwareProfiles))
		is.Equal(12, len(defs.Distributions))

		for k, l := range fixtureDefinitions.Distributions {
			is.Equal(l, defs.Distributions[k])
		}
		for k, l := range fixtureDefinitions.StorageGrades {
			is.Equal(l, defs.StorageGrades[k])
		}
		for k, l := range fixtureDefinitions.HardwareProfiles {
			is.Equal(l, defs.HardwareProfiles[k])
		}
		for k, l := range fixtureDefinitions.ZoneNames {
			is.Equal(l, defs.ZoneNames[k])
		}
	})

}

var fixtureDefinitions = lib.Definitions{
	DistributionDescriptions: map[string]string{
		"centos5":     "CentOS 5 (64 bit)",
		"centos6":     "CentOS 6 (64 bit)",
		"centos7":     "CentOS 7 (64 bit)",
		"jessie":      "Debian 8 (64 bit)",
		"precise":     "Ubuntu 12.04 (64 bit)",
		"symbiosis":   "Symbiosis 7 (64 bit)",
		"trusty":      "Ubuntu 14.04 (64 bit)",
		"utopic":      "Ubuntu 14.10 (64 bit)",
		"vivid":       "Ubuntu 15.04 (64 bit)",
		"wheezy":      "Debian 7 (64 bit)",
		"winstd2012":  "Windows Server Standard 2012 (64 bit)",
		"winweb2k8r2": "Windows Web Server 2008R2 (64 bit)",
	},
	StorageGradeDescriptions: map[string]string{
		"archive": "Archive storage",
		"sata":    "Standard SSD",
		"ssd":     "Premium SSD",
	},
	HardwareProfiles: []string{
		"compatibility2013",
		"virtio2013",
	},
	ZoneNames: []string{
		"manchester",
		"york",
	},
}

const fixtureDefinitionsJSON = `[
    {
        "id": "distributions",
        "data": [
            "centos5",
            "centos6",
            "winstd2012",
            "centos7",
            "jessie",
            "precise",
            "symbiosis",
            "trusty",
            "utopic",
            "vivid",
            "wheezy",
            "winweb2k8r2"
        ]
    },
    {
        "id": "storage_grades",
        "data": [
            "sata",
            "ssd",
            "archive"
        ]
    },
    {
        "id": "zone_names",
        "data": [
            "manchester",
            "york"
        ]
    },
    {
        "id": "distribution_descriptions",
        "data": {
            "centos5": "CentOS 5 (64 bit)",
            "symbiosis": "Symbiosis 7 (64 bit)",
            "centos6": "CentOS 6 (64 bit)",
            "centos7": "CentOS 7 (64 bit)",
            "jessie": "Debian 8 (64 bit)",
            "precise": "Ubuntu 12.04 (64 bit)",
            "trusty": "Ubuntu 14.04 (64 bit)",
            "utopic": "Ubuntu 14.10 (64 bit)",
            "vivid": "Ubuntu 15.04 (64 bit)",
            "wheezy": "Debian 7 (64 bit)",
            "winstd2012": "Windows Server Standard 2012 (64 bit)",
            "winweb2k8r2": "Windows Web Server 2008R2 (64 bit)"
        }
    },
    {
        "id": "storage_grade_descriptions",
        "data": {
            "ssd": "Premium SSD",
            "archive": "Archive storage",
            "sata": "Standard SSD"
        }
    },
    {
        "id": "hardware_profiles",
        "data": [
            "virtio2013",
            "compatibility2013"
        ]
    },
    {
        "id": "keymaps",
        "data": [
            "ar",
            "de-ch",
            "es",
            "fo",
            "fr-ca",
            "hu",
            "ja",
            "mk",
            "no",
            "pt-br",
            "sv",
            "da",
            "en-gb",
            "et",
            "fr",
            "fr-ch",
            "is",
            "lt",
            "nl",
            "pl",
            "ru",
            "th",
            "de",
            "en-us",
            "fi",
            "fr-be",
            "hr",
            "it",
            "lv",
            "nl-be",
            "pt",
            "sl",
            "tr"
        ]
    },
    {
        "id": "sendkeys",
        "data": [
            "shift",
            "shift_r",
            "alt",
            "alt_r",
            "altgr",
            "altgr_r",
            "ctrl",
            "ctrl_r",
            "menu",
            "esc",
            "1",
            "2",
            "3",
            "4",
            "5",
            "6",
            "7",
            "8",
            "9",
            "0",
            "minus",
            "equal",
            "backspace",
            "tab",
            "q",
            "w",
            "e",
            "r",
            "t",
            "y",
            "u",
            "i",
            "o",
            "p",
            "ret",
            "a",
            "s",
            "d",
            "f",
            "g",
            "h",
            "j",
            "k",
            "l",
            "z",
            "x",
            "c",
            "v",
            "b",
            "n",
            "m",
            "comma",
            "dot",
            "slash",
            "asterisk",
            "spc",
            "caps_lock",
            "f1",
            "f2",
            "f3",
            "f4",
            "f5",
            "f6",
            "f7",
            "f8",
            "f9",
            "f10",
            "f11",
            "f12",
            "num_lock",
            "scroll_lock",
            "kp_divide",
            "kp_multiply",
            "kp_subtract",
            "kp_add",
            "kp_enter",
            "kp_decimal",
            "sysrq",
            "kp_0",
            "kp_1",
            "kp_2",
            "kp_3",
            "kp_4",
            "kp_5",
            "kp_6",
            "kp_7",
            "kp_8",
            "kp_9",
            "print",
            "home",
            "pgup",
            "pgdn",
            "end",
            "left",
            "up",
            "down",
            "right",
            "insert",
            "delete"
        ]
    }
]`
