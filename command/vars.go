package command

import (
	"github.com/rs/zerolog"
	"github.com/urfave/cli"
)

var commonFlags = []cli.Flag{
	cli.BoolFlag{
		Name:  "debug",
		Usage: "print debug logs",
	},
}

func setLogLevel(ctx *cli.Context) {
	if ctx.Bool("debug") {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}

var (
	// VersionNumber is version string of PuGo
	VersionNumber = "2.3.3"
)
