package main

import (
	"flag"
	"github.com/VictoriaMetrics/VictoriaMetrics/app/vminsert/common"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/logger"
)

func main() {
	flag.Parse()
	logger.Init()
	if err := common.ConfigureInsert(); err != nil {
		logger.Fatalf("failed to configure insert: %v", err)
	}
	logger.Infof("Starting vminsert...")
}
