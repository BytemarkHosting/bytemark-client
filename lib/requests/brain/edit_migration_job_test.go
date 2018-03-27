package brain_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	brainMethods "github.com/BytemarkHosting/bytemark-client/lib/requests/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil/assert"
)

func TestEditMigrationJob(t *testing.T) {
	tests := []struct {
		id            int
		modifications brain.MigrationJobModification
		expected      map[string]interface{}
		shouldErr     bool
	}{
		{
			id: 1,
			modifications: brain.MigrationJobModification{
				Cancel: brain.MigrationJobLocations{
					Discs: []json.Number{"disc.sata-1.8912"},
					Pools: []json.Number{"t1-archive1"},
					Tails: []json.Number{"tail2"}},
				Options: brain.MigrationJobOptions{
					Priority: 10,
				},
			},
			expected: map[string]interface{}{
				"Cancel": map[string]interface{}{
					"Discs": "disc.sata-1.8912",
					"Tails": "t1-archive1",
					"Pools": "tail2",
				},
				"Options": map[string]interface{}{
					"Priority": 10,
				},
			},
		},
	}
	for i, test := range tests {
		testName := testutil.Name(i)
		rts := testutil.RequestTestSpec{
			Method:        "PUT",
			Endpoint:      lib.BrainEndpoint,
			URL:           fmt.Sprintf("/admin/migration_jobs/%s", test.id),
			AssertRequest: assert.BodyUnmarshalEqual(test.expected),
		}
		rts.Run(t, testName, true, func(client lib.Client) {
			err := brainMethods.EditMigrationJob(client, test.id, test.modifications)
			if test.shouldErr {
				assert.NotEqual(t, testName, nil, err)
			} else {
				assert.Equal(t, testName, nil, err)
			}
		})
	}
}
