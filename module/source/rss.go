package source

import (
	"os"
	"path/filepath"
	"time"

	"github.com/fuxiaohei/pugov1/module/config"
	"github.com/gorilla/feeds"
	"github.com/rs/zerolog/log"
)

func renderRSS(ctx *compileContext) {
	now := time.Now()
	feed := &feeds.Feed{
		Title:       ctx.Config.Meta.Title,
		Link:        &feeds.Link{Href: ctx.Config.Meta.Root},
		Description: ctx.Config.Meta.Desc,
		Created:     now,
	}
	owner := config.Owner(ctx.Config)
	if owner != nil {
		feed.Author = &feeds.Author{
			Name:  owner.Nick,
			Email: owner.Email,
		}
	}
	var item *feeds.Item
	for _, p := range ctx.Source.Posts {
		item = &feeds.Item{
			Title:       p.Title,
			Link:        &feeds.Link{Href: fullURL(ctx.Config.Meta.Root, p.URL)},
			Description: string(p.ContentBytes),
			Created:     p.Created(),
			Updated:     p.Updated(),
		}
		if p.Author != nil {
			item.Author = &feeds.Author{
				Name:  p.Author.Nick,
				Email: p.Author.Email,
			}
		}
		feed.Items = append(feed.Items, item)
	}

	dstFile := filepath.Join(ctx.OutputDir, "feed.xml")
	f, err := os.OpenFile(dstFile, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		log.Warn().Err(err).Msg("render-rss-error")
		return
	}
	defer f.Close()
	if err = feed.WriteRss(f); err != nil {
		log.Warn().Err(err).Msg("render-rss-error")
		return
	}
	log.Debug().Str("dest", dstFile).Msg("render-rss-ok")
	ctx.Source.RenderedFiles = append(ctx.Source.RenderedFiles, dstFile)
}
