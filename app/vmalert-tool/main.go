package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/VictoriaMetrics/VictoriaMetrics/lib/buildinfo"
	"github.com/urfave/cli/v2"
)

func main() {
	start := time.Now()
	app := &cli.App{
		Name:    "vmalert-tool",
		Usage:   "VMAlert command-line tool",
		Version: buildinfo.Version,
		Commands: []*cli.Command{
			{
				Name:  "unittest",
				Usage: "Run unittests for alerting and recording rules.",
				Action: func(c *cli.Context) error {
					fmt.Println("Running unittests...")
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Total time: %v", time.Since(start))
}
