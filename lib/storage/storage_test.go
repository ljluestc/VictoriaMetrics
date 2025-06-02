package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestEnforceRetention(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "vmstorage-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	partitionsDir := filepath.Join(tmpDir, "data", "big")
	if err := os.MkdirAll(partitionsDir, 0755); err != nil {
		t.Fatalf("Failed to create partitions dir: %v", err)
	}
	for _, month := range []string{"202501", "202502"} {
		if err := os.Mkdir(filepath.Join(partitionsDir, month), 0755); err != nil {
			t.Fatalf("Failed to create partition %s: %v", month, err)
		}
	}

	config := &Config{
		StorageDataPath: tmpDir,
		RetentionPeriod: "1", // 1 month.
		KeepOldData:     false,
	}
	if err := enforceRetention(config); err != nil {
		t.Fatalf("Retention failed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(partitionsDir, "202501")); !os.IsNotExist(err) {
		t.Errorf("Partition 202501 was not deleted")
	}
	if _, err := os.Stat(filepath.Join(partitionsDir, "202502")); os.IsNotExist(err) {
		t.Errorf("Partition 202502 was deleted unexpectedly")
	}

	config.KeepOldData = true
	if err := os.Mkdir(filepath.Join(partitionsDir, "202501"), 0755); err != nil {
		t.Fatalf("Failed to recreate partition: %v", err)
	}
	if err := enforceRetention(config); err != nil {
		t.Fatalf("Retention failed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(partitionsDir, "202501")); os.IsNotExist(err) {
		t.Errorf("Partition 202501 was deleted despite keepOldData=true")
	}
}
