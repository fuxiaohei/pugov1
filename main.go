package main

//go:generate ./pugov1 pack

import (
	"os"

	_ "github.com/fuxiaohei/pugov1/asset"
	"github.com/fuxiaohei/pugov1/command"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.TimeFieldFormat = "01/02 15:04:05"
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:     os.Stderr,
		NoColor: false,
	})
}

func main() {

	app := cli.NewApp()
	app.Name = "PuGo"
	app.Usage = "a simple static site generator"
	app.Description = app.Usage
	app.Commands = []cli.Command{
		command.Init,
		command.Build,
		command.Server,
		command.Pack,
		command.Version,
	}
	app.RunAndExitOnError()
}
