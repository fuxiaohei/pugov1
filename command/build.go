package command

import (
	"path/filepath"
	"time"

	"github.com/fuxiaohei/pugov1/module/config"
	"github.com/fuxiaohei/pugov1/module/source"
	"github.com/fuxiaohei/pugov1/module/theme"
	"github.com/fuxiaohei/pugov1/module/watch"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli"
)

// Build is command of 'build'
var Build = cli.Command{
	Name:  "build",
	Usage: "build source and theme to website pages",
	Flags: append(commonFlags, cli.BoolFlag{
		Name:  "clean",
		Usage: "clean not-rendered file",
	}, cli.BoolFlag{
		Name:  "watch",
		Usage: "watch files changes to rebuild",
	}),
	Action: func(ctx *cli.Context) error {
		setLogLevel(ctx)
		if ctx.Bool("watch") {
			watchOnce([]string{sourceDir, themeDir}, ctx)
			return nil
		}
		return buildOnce(ctx)
	},
}

func watchOnce(dirs []string, ctx *cli.Context) {
	watch.Watch(dirs, func(_ []*watch.Event) {
		if err := buildOnce(ctx); err != nil {
			log.Warn().Err(err).Msg("build-error")
		}
	})
}

func buildOnce(ctx *cli.Context) error {
	st := time.Now()
	if err := buildFunc(ctx); err != nil {
		return err
	}
	log.Info().
		Dur("duration", time.Since(st)).
		Msg("build-ok")
	return nil
}

var (
	sourceDir = "source"
	themeDir  = "theme"

	postDir  = "post"
	pageDir  = "page"
	langDir  = "lang"
	adminDir = "admin"

	outputDir = "_dest"
)

func buildFunc(ctx *cli.Context) error {
	cfg, err := config.Read()
	if err != nil {
		return err
	}
	log.Info().Str("file", cfg.SrcFile).Msg("read-config")

	th, err := theme.Read(themeDir, "0.11.0")
	if err != nil {
		return err
	}
	log.Info().Str("directory", th.Dir()).Msg("read-theme")

	s, err := source.Read(
		filepath.Join(sourceDir, postDir),
		filepath.Join(sourceDir, pageDir),
		filepath.Join(sourceDir, langDir),
	)
	if err != nil {
		return err
	}
	if err = source.Parse(s, false); err != nil {
		return err
	}
	log.Info().
		Int("posts", len(s.Posts)).
		Int("pages", len(s.Pages)).
		Int("lists", len(s.Lists)).
		Int("tags", len(s.TagsLists)).
		Int("langs", len(s.I18ns)).
		Msg("read-contents")

	if err = source.Compile(s, cfg, th, outputDir); err != nil {
		return err
	}
	count, skip, err := source.Copy(s, outputDir)
	if err != nil {
		return err
	}
	log.Info().
		Int("files", count).
		Int("skip", skip).
		Msg("copy-files")
	if ctx.Bool("clean") {
		count, err = source.Cleanup(s, outputDir)
		if err != nil {
			return err
		}
		log.Info().
			Int("files", count).
			Msg("clean-files")
	}
	return nil
}
