package command

import (
	"github.com/fuxiaohei/pugov1/asset"
	"github.com/fuxiaohei/pugov1/module/packer"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli"
)

// Init is command of 'init'
var Init = cli.Command{
	Name:  "init",
	Usage: "init default site",
	Flags: commonFlags,
	Action: func(ctx *cli.Context) error {
		setLogLevel(ctx)
		return unpackOnce(ctx)
	},
}

func unpackOnce(ctx *cli.Context) error {
	log.Info().Int("files", len(asset.Files)).Msg("load-asset")
	return packer.Unasset(asset.Files)
}
