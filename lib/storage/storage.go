package storage

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

// Fix enforceRetention function
func enforceRetention(config *Config) error {
	retention, err := parseDuration(config.RetentionPeriod)
	if err != nil {
		return fmt.Errorf("failed to parse retention period: %v", err)
	}

	cutoffTime := time.Now().Add(-retention)
	partitionsDir := filepath.Join(config.StorageDataPath, "data", "big") // Example path.

	return filepath.Walk(partitionsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return nil
		}

		partitionTime, err := time.Parse("200601", info.Name())
		if err != nil {
			return nil // Skip invalid partition names.
		}

		if partitionTime.Before(cutoffTime) {
			if config.KeepOldData {
				log.Printf("Retaining partition %s due to -keepOldData flag", info.Name())
				return nil
			}
			log.Printf("Deleting partition %s older than retention period %s", info.Name(), config.RetentionPeriod)
			return os.RemoveAll(path) // Replace fs.RemoveDir with os.RemoveAll
		}
		return nil
	})
}
