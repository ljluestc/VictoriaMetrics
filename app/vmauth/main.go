package main

import (
	"flag"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/logger"
)

func main() {
	flag.Parse()
	logger.Init()
	initAuthConfig()
	logger.Infof("Starting vmauth...")
}
