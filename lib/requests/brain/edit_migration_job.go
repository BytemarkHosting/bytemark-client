package brain

import (
	"strconv"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
)

// EditMigrationJob allows you to cancel individual or multiples discs, pools, or tails and change the priority for a job, given its ID
func EditMigrationJob(client lib.Client, id int, migrationEdit brain.MigrationJobModification) (err error) {
	r, err := client.BuildRequest("PUT", lib.BrainEndpoint, "/admin/migration_jobs/%s", strconv.Itoa(id))
	if err != nil {
		return
	}

	_, _, err = r.MarshalAndRun(migrationEdit, nil)
	return
}

// CancelMigrationJob cancels all migrations on a job, given its ID
func CancelMigrationJob(client lib.Client, id int) (err error) {
	r, err := client.BuildRequest("PUT", lib.BrainEndpoint, "/admin/migration_jobs/%s", strconv.Itoa(id))
	if err != nil {
		return
	}

	cancel := map[string]interface{}{
		"cancel": map[string]interface{}{
			"all": true,
		},
	}

	_, _, err = r.MarshalAndRun(cancel, nil)
	return
}
