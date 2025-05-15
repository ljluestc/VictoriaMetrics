package vmstorage

import (
	"flag"
)

// Define variables
var (
	keepOldData  = flag.Bool("keepOldData", false, "If set, old data won't be removed on startup when a smaller retention period is provided")
	oldRetention = flag.Int("oldRetention", 0, "Old retention period")
	newRetention = flag.Int("newRetention", 0, "New retention period")
)

func ExampleFunction() {
	if *oldRetention > *newRetention && !*keepOldData {
		// logic to delete old data
	}
}

func updateRetention() {
	// Logic to apply retention
}
