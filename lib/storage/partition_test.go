package storage

import (
	"testing"
	"time"

	"github.com/VictoriaMetrics/VictoriaMetrics/lib/fs"
)

func TestRemoveStalePartsWithKeepOldData(t *testing.T) {
	dataPath := t.TempDir()
	retentionPeriod := "1d"
	keepOldData := true

	if err := Init(dataPath, retentionPeriod, keepOldData); err != nil {
		t.Fatalf("failed to init storage: %v", err)
	}

	timestamp := time.Now().Add(-48 * time.Hour).UnixMilli() // 2 days ago
	pt := mustCreatePartition(timestamp, dataPath+"/small", dataPath+"/big", &Storage{
		retentionMsecs: int64(24 * time.Hour / time.Millisecond),
		keepOldData:    true,
	})

	rows := []rawRow{{Timestamp: timestamp, PrecisionBits: 32}}
	pt.AddRows(rows)
	pt.flushRowssToInmemoryParts([][]rawRow{rows})

	pt.removeStaleParts()

	pt.partsLock.Lock()
	if len(pt.inmemoryParts) == 0 {
		t.Fatalf("expected part to be preserved with keepOldData=true")
	}
	pt.partsLock.Unlock()

	pt.s.keepOldData = false
	pt.removeStaleParts()

	pt.partsLock.Lock()
	if len(pt.inmemoryParts) != 0 {
		t.Fatalf("expected part to be deleted with keepOldData=false")
	}
	pt.partsLock.Unlock()

	pt.MustClose()
	fs.MustRemoveAll(dataPath)
}
