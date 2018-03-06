package command

import (
	"fmt"
	"runtime"

	"github.com/fuxiaohei/pugov1/asset"
	"github.com/urfave/cli"
)

// Version is command of 'version'
var Version = cli.Command{
	Name:  "version",
	Usage: "print version",
	Flags: commonFlags,
	Action: func(ctx *cli.Context) error {
		setLogLevel(ctx)
		return printVersion(ctx)
	},
}

func printVersion(ctx *cli.Context) error {
	if ctx.Bool("debug") {
		fmt.Printf("PuGo:\t%s\n", VersionNumber)
		fmt.Printf("Go:\t%s\n", runtime.Version())
		fmt.Printf("Assets:\t%s\n", asset.FilesTime)
		return nil
	}
	fmt.Println(VersionNumber)
	return nil
}
