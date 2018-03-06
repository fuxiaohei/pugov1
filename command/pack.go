package command

import (
	"github.com/fuxiaohei/pugov1/module/packer"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli"
)

// Pack is command of 'build'
var Pack = cli.Command{
	Name:  "pack",
	Usage: "pack source code to go file",
	Flags: commonFlags,
	Action: func(ctx *cli.Context) error {
		setLogLevel(ctx)
		return packOnce(ctx)
	},
}

var (
	assetFile = "asset/asset.go"
)

func packOnce(ctx *cli.Context) error {
	filesData := packer.Files([]string{"config.toml"}, []string{
		sourceDir,
		themeDir,
	})
	n, err := packer.Asset(assetFile, filesData)
	if err != nil {
		return err
	}
	log.Info().Str("file", assetFile).Int("bytes", n/1024).Msg("pack-to-file")
	return nil
}
