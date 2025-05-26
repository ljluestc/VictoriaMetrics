package vmselect

import (
	"flag"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/logger"
)

func main() {
	flag.Parse()
	logger.Init()
	logger.Infof("Starting vmselect...")
}
