package main

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"github.com/go-logr/logr"
	"github.com/urfave/cli/v2"
)

var (
	// these variables are populated by Goreleaser when releasing
	version = "unknown"
	commit  = "-dirty-"
	date    = time.Now().Format("2006-01-02")

	appName     = "provider-cloudscale"
	appLongName = "Crossplane provider that deploys resources on cloudscale.ch"
)

func main() {
	app := newApp()
	err := app.Run(os.Args)
	// If required flags aren't set, it will return with error before we could set up logging
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func newApp() *cli.App {
	logInstance := &atomic.Value{}
	logInstance.Store(logr.Discard())
	app := &cli.App{
		Name:    appName,
		Usage:   appLongName,
		Version: fmt.Sprintf("%s, revision=%s, date=%s", version, commit, date),

		Before: setupLogging,
		Flags: []cli.Flag{
			newLogLevelFlag(),
			newLogFormatFlag(),
		},
		Commands: []*cli.Command{
			newOperatorCommand(),
		},
	}
	return app
}
