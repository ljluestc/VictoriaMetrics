package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCleanOldData(t *testing.T) {
	dataPath := filepath.Join(os.TempDir(), "vm-test-storage")
	defer os.RemoveAll(dataPath)

	partitions := []string{"202301", "202505", "invalid"}
	for _, p := range partitions {
		partitionPath := filepath.Join(dataPath, "data", p)
		if err := os.MkdirAll(partitionPath, 0755); err != nil {
			t.Fatalf("cannot create partition %q: %s", partitionPath, err)
		}
	}

	t.Run("WithoutKeepOldData", func(t *testing.T) {
		if err := cleanOldData(dataPath, 30*24*time.Hour, false); err != nil {
			t.Fatalf("cleanOldData failed: %s", err)
		}
		if _, err := os.Stat(filepath.Join(dataPath, "data", "202301")); !os.IsNotExist(err) {
			t.Errorf("old partition 202301 was not deleted")
		}
	})

	t.Run("WithKeepOldData", func(t *testing.T) {
		os.MkdirAll(filepath.Join(dataPath, "data", "202301"), 0755)
		if err := cleanOldData(dataPath, 30*24*time.Hour, true); err != nil {
			t.Fatalf("cleanOldData failed: %s", err)
		}
		if _, err := os.Stat(filepath.Join(dataPath, "data", "202301")); os.IsNotExist(err) {
			t.Errorf("old partition 202301 was deleted with -keepOldData")
		}
	})
}
