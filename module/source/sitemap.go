package source

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

func fullURL(domain string, urlStr string) string {
	return strings.TrimSuffix(domain, "/") + "/" + strings.TrimPrefix(urlStr, "/")
}

func renderSiteMap(ctx *compileContext) {
	now := time.Now()
	var buf bytes.Buffer
	buf.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	buf.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`)
	buf.WriteString("<url>")
	fmt.Fprintf(&buf, "<loc>%s</loc>", ctx.Config.Meta.Root)
	fmt.Fprintf(&buf, "<lastmod>%s</lastmod>", now.Format(time.RFC3339))
	buf.WriteString("<changefreq>daily</changefreq>")
	buf.WriteString("<priority>1.0</priority>")
	buf.WriteString("</url>")

	buf.WriteString("<url>")
	fmt.Fprintf(&buf, "<loc>%s</loc>", fullURL(ctx.Config.Meta.Root, "archive.html"))
	fmt.Fprintf(&buf, "<lastmod>%s</lastmod>", now.Format(time.RFC3339))
	buf.WriteString("<changefreq>daily</changefreq>")
	buf.WriteString("<priority>0.6</priority>")
	buf.WriteString("</url>")

	for _, p := range ctx.Source.Pages {
		buf.WriteString("<url>")
		fmt.Fprintf(&buf, "<loc>%s</loc>", fullURL(ctx.Config.Meta.Root, p.URL))
		fmt.Fprintf(&buf, "<lastmod>%s</lastmod>", p.Created().Format(time.RFC3339))
		buf.WriteString("<changefreq>weekly</changefreq>")
		buf.WriteString("<priority>0.5</priority>")
		buf.WriteString("</url>")
	}

	for _, p := range ctx.Source.Posts {
		buf.WriteString("<url>")
		fmt.Fprintf(&buf, "<loc>%s</loc>", fullURL(ctx.Config.Meta.Root, p.URL))
		fmt.Fprintf(&buf, "<lastmod>%s</lastmod>", p.Created().Format(time.RFC3339))
		buf.WriteString("<changefreq>daily</changefreq>")
		buf.WriteString("<priority>0.6</priority>")
		buf.WriteString("</url>")
	}

	for i := 1; i <= ctx.Source.ListSize; i++ {
		buf.WriteString("<url>")
		fmt.Fprintf(&buf, "<loc>%s</loc>", fullURL(ctx.Config.Meta.Root, fmt.Sprintf("post/%d.html", i)))
		fmt.Fprintf(&buf, "<lastmod>%s</lastmod>", now.Format(time.RFC3339))
		buf.WriteString("<changefreq>daily</changefreq>")
		buf.WriteString("<priority>0.6</priority>")
		buf.WriteString("</url>")
	}

	for _, list := range ctx.Source.TagsLists {
		buf.WriteString("<url>")
		fmt.Fprintf(&buf, "<loc>%s</loc>", fullURL(ctx.Config.Meta.Root, list.Tag.URL))
		fmt.Fprintf(&buf, "<lastmod>%s</lastmod>", now.Format(time.RFC3339))
		buf.WriteString("<changefreq>weekly</changefreq>")
		buf.WriteString("<priority>0.5</priority>")
		buf.WriteString("</url>")
	}

	buf.WriteString("</urlset>")
	dstFile := filepath.Join(ctx.OutputDir, "sitemap.xml")
	os.MkdirAll(path.Dir(dstFile), os.ModePerm)
	if err := ioutil.WriteFile(dstFile, buf.Bytes(), os.ModePerm); err != nil {
		log.Warn().Err(err).Msg("render-sitemap-error")
		return
	}
	log.Debug().
		Str("dest", dstFile).Msg("render-sitemap-ok")
	ctx.Source.RenderedFiles = append(ctx.Source.RenderedFiles, dstFile)
}
