package main

import (
	"flag"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/logger"
)

func main() {
	logger.Init()
	flag.Parse()
	logger.Infof("Starting vmagent...")
}
