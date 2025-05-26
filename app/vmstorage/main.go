package vmstorage

import (
	"flag"

	"github.com/VictoriaMetrics/VictoriaMetrics/lib/logger"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/storage"
)

var (
	retentionPeriod = flag.String("retentionPeriod", "1", "Retention period in months (e.g., 1, 3d, 2w, 100y; minimum 1d)")
	storageDataPath = flag.String("storageDataPath", "victoria-metrics-data", "Path to storage data")
	keepOldData     = flag.Bool("keepOldData", false, "Prevent deletion of data outside retention period when reducing -retentionPeriod")
)

func InitStorage() error {
	flag.Parse()
	logger.Init()

	logger.Infof("Starting vmstorage with retentionPeriod=%s and keepOldData=%v", *retentionPeriod, *keepOldData)

	if err := storage.Init(*storageDataPath, *retentionPeriod, *keepOldData); err != nil {
		return err
	}
	return nil
}
