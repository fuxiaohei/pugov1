package source

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/fuxiaohei/pugov1/module/config"
	"github.com/fuxiaohei/pugov1/module/i18n"
	"github.com/fuxiaohei/pugov1/module/object"
	"github.com/fuxiaohei/pugov1/module/theme"
	"github.com/rs/zerolog/log"
)

// Compile compile source to html files
func Compile(s *object.Source, cfg *object.Config, th *theme.Theme, outputDir string) error {
	os.MkdirAll(outputDir, os.ModePerm)
	ctx := &compileContext{
		Source:    s,
		Config:    cfg,
		Theme:     th,
		OutputDir: outputDir,
	}
	staticFiles, err := th.StaticFiles()
	if err != nil {
		return err
	}
	s.CopyFiles = append(s.CopyFiles, staticFiles...)

	ctx.renderPosts()
	ctx.renderLists()
	ctx.renderTagLists()
	ctx.renderArchive()
	ctx.renderIndex()
	ctx.renderPages()

	renderSiteMap(ctx)
	renderRSS(ctx)

	return nil
}

type compileContext struct {
	Source    *object.Source
	Config    *object.Config
	Theme     *theme.Theme
	OutputDir string
}

func (ctx *compileContext) viewData() map[string]interface{} {
	cfg := ctx.Config
	m := map[string]interface{}{
		"Source":    ctx.Source,
		"Nav":       cfg.Navs,
		"Meta":      cfg.Meta,
		"Title":     cfg.Meta.Title + " - " + cfg.Meta.Subtitle,
		"Desc":      cfg.Meta.Desc,
		"Comment":   cfg.Comment,
		"Owner":     config.Owner(cfg), //ctx.Source.Owner,
		"Analytics": cfg.Analytics,
		"Lang":      cfg.Meta.Lang,
		"Hover":     "",
		"Root":      strings.TrimRight(cfg.Meta.Root, "/"),
	}
	in := ctx.Source.I18ns[cfg.Meta.Lang]
	if in == nil {
		in = i18n.Empty()
	}
	m["I18n"] = in
	return m
}

func (ctx *compileContext) renderPosts() {
	for _, p := range ctx.Source.Posts {
		dstFile := filepath.Join(ctx.OutputDir, p.URL)
		p.OutputFile = dstFile

		buf := bytes.NewBuffer(nil)
		viewData := ctx.viewData()
		viewData["Post"] = p

		if err := ctx.Theme.Execute(buf, "post.html", viewData); err != nil {
			log.Warn().Str("file", p.SourceFile).Err(err).Msg("render-post-error")
			continue
		}
		dir := filepath.Dir(dstFile)
		os.MkdirAll(dir, os.ModePerm)
		if err := ioutil.WriteFile(dstFile, buf.Bytes(), 0644); err != nil {
			log.Warn().Str("file", p.SourceFile).Err(err).Msg("render-post-error")
			continue
		}
		log.Debug().
			Str("file", p.SourceFile).
			Str("dest", dstFile).
			Msg("render-post-ok")
		ctx.Source.RenderedFiles = append(ctx.Source.RenderedFiles, dstFile)
	}
}

func (ctx *compileContext) renderLists() {
	for _, list := range ctx.Source.Lists {
		permaURL := fmt.Sprintf("post/%d.html", list.Pager.Current)
		dstFile := filepath.Join(ctx.OutputDir, permaURL)

		list.Pager.SetLayout("post/%d.html")

		buf := bytes.NewBuffer(nil)
		viewData := ctx.viewData()
		viewData["Posts"] = list.Posts
		viewData["Pager"] = list.Pager
		if err := ctx.Theme.Execute(buf, "posts.html", viewData); err != nil {
			log.Warn().Int("page", list.Pager.Current).Err(err).Msg("render-list-error")
			continue
		}
		dir := filepath.Dir(dstFile)
		os.MkdirAll(dir, os.ModePerm)
		if err := ioutil.WriteFile(dstFile, buf.Bytes(), 0644); err != nil {
			log.Warn().Int("page", list.Pager.Current).Err(err).Msg("render-list-error")
			continue
		}
		log.Debug().
			Str("dest", dstFile).
			Msg("render-list-ok")
		ctx.Source.RenderedFiles = append(ctx.Source.RenderedFiles, dstFile)
	}
}

func (ctx *compileContext) renderTagLists() {
	for _, list := range ctx.Source.TagsLists {
		permaURL := fmt.Sprintf("/tags/%s.html", list.Tag.Name)
		dstFile := filepath.Join(ctx.OutputDir, permaURL)

		buf := bytes.NewBuffer(nil)
		viewData := ctx.viewData()
		viewData["Archives"] = list.Archives
		viewData["Tag"] = list.Tag

		if err := ctx.Theme.Execute(buf, "archive.html", viewData); err != nil {
			log.Warn().Str("tag", list.Tag.Name).Err(err).Msg("render-tag-list-error")
			continue
		}
		dir := filepath.Dir(dstFile)
		os.MkdirAll(dir, os.ModePerm)
		if err := ioutil.WriteFile(dstFile, buf.Bytes(), 0644); err != nil {
			log.Warn().Str("tag", list.Tag.Name).Err(err).Msg("render-tag-list-error")
			continue
		}
		log.Debug().
			Str("tag", list.Tag.Name).
			Str("dest", dstFile).
			Msg("render-tag-list-ok")
		ctx.Source.RenderedFiles = append(ctx.Source.RenderedFiles, dstFile)
	}
}

func (ctx *compileContext) renderArchive() {
	dstFile := filepath.Join(ctx.OutputDir, "archive.html")
	buf := bytes.NewBuffer(nil)
	viewData := ctx.viewData()
	viewData["Archives"] = ctx.Source.Archives
	viewData["Hover"] = "archive"
	if err := ctx.Theme.Execute(buf, "archive.html", viewData); err != nil {
		log.Warn().Err(err).Msg("render-archive-error")
		return
	}
	dir := filepath.Dir(dstFile)
	os.MkdirAll(dir, os.ModePerm)
	if err := ioutil.WriteFile(dstFile, buf.Bytes(), 0644); err != nil {
		log.Warn().Err(err).Msg("render-archive-error")
		return
	}
	log.Debug().
		Str("dest", dstFile).
		Msg("render-archive-ok")
	ctx.Source.RenderedFiles = append(ctx.Source.RenderedFiles, dstFile)
}

func (ctx *compileContext) renderIndex() {
	dstFile := filepath.Join(ctx.OutputDir, "index.html")
	buf := bytes.NewBuffer(nil)
	tpl := ctx.Theme.Template("index.html")
	if tpl == nil {
		tpl = ctx.Theme.Template("posts.html")
	}
	viewData := ctx.viewData()
	viewData["Posts"] = ctx.Source.Lists[0].Posts
	viewData["Pager"] = ctx.Source.Lists[0].Pager
	viewData["Hover"] = "home"
	if err := tpl.Execute(buf, viewData); err != nil {
		log.Warn().Err(err).Msg("render-index-error")
		return
	}
	dir := filepath.Dir(dstFile)
	os.MkdirAll(dir, os.ModePerm)
	if err := ioutil.WriteFile(dstFile, buf.Bytes(), 0644); err != nil {
		log.Warn().Err(err).Msg("render-index-error")
		return
	}
	log.Debug().
		Str("dest", dstFile).
		Msg("render-index-ok")
	ctx.Source.RenderedFiles = append(ctx.Source.RenderedFiles, dstFile)
}

func (ctx *compileContext) renderPages() {
	for _, p := range ctx.Source.Pages {
		dstFile := filepath.Join(ctx.OutputDir, p.URL)
		p.OutputFile = dstFile

		buf := bytes.NewBuffer(nil)
		viewData := ctx.viewData()
		viewData["Page"] = p
		viewData["Hover"] = p.NavHover
		if p.Lang != "" {
			in := ctx.Source.I18ns[p.Lang]
			if in != nil {
				viewData["I18n"] = in
			} else {
				log.Warn().Str("file", p.SourceFile).Str("i18n", p.Lang).Msg("render-page-i18n-missing")
			}
		}

		if err := ctx.Theme.Execute(buf, "page.html", viewData); err != nil {
			log.Warn().Str("file", p.SourceFile).Err(err).Msg("render-page-error")
			continue
		}
		dir := filepath.Dir(dstFile)
		os.MkdirAll(dir, os.ModePerm)
		if err := ioutil.WriteFile(dstFile, buf.Bytes(), 0644); err != nil {
			log.Warn().Str("file", p.SourceFile).Err(err).Msg("render-page-error")
			continue
		}
		log.Debug().
			Str("file", p.SourceFile).
			Str("dest", dstFile).
			Msg("render-page-ok")
		ctx.Source.RenderedFiles = append(ctx.Source.RenderedFiles, dstFile)
	}
}
