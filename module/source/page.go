package source

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/fuxiaohei/pugov1/module/object"
	"github.com/rs/zerolog/log"
	blackfriday "gopkg.in/russross/blackfriday.v2"
)

func parsePages(s *object.Source, withDraft bool) ([]*object.Page, error) {
	var pages []*object.Page
	for _, item := range s.PageFiles {
		relPath, _ := filepath.Rel(s.PageDir, item.File)

		if item.OpType == object.SourceOpFile {
			item.SrcFile = item.File
			item.File = relPath
			s.CopyFiles = append(s.CopyFiles, item)
			continue
		}
		page, err := parseOnePage(item.File, relPath, item.Info)
		if err != nil {
			log.Warn().Str("src", item.File).Err(err).Msg("read-page-error")
			continue
		}
		page.SourceRelpath, _ = filepath.Rel(s.PageDir, page.SourceFile)
		if page.Draft {
			log.Warn().Str("src", item.File).Msg("page-draft")
			if !withDraft {
				continue
			}
		}
		pages = append(pages, page)
	}
	return pages, nil
}

const (
	pageSeparator = "```"
	pageMetaTOML  = "toml"
)

var (
	// ErrorPageSectionsInvalid means post cant not seperate to current sections for parsing
	ErrorPageSectionsInvalid = errors.New("page sections are invalid")
	// ErrorPageMetaUnknownType means post metadata cant be unmarshaled with proper type
	ErrorPageMetaUnknownType = errors.New("page metadata is unknown type")
)

func parseOnePage(file string, relPath string, fileInfo os.FileInfo) (*object.Page, error) {
	fData, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	sections := bytes.SplitN(fData, []byte(pageSeparator), 3)
	if len(sections) < 3 {
		return nil, ErrorPageSectionsInvalid
	}
	p := new(object.Page)
	if err := parsePageMeta(sections[1], p); err != nil {
		return nil, err
	}
	p.RawBytes = bytes.TrimSpace(sections[2])
	p.SourceFile = file
	p.URL = strings.TrimSuffix(relPath, filepath.Ext(relPath)) + ".html"
	if p.Template == "" {
		p.Template = "page.html"
	}
	return p, fillOnePage(p, fileInfo)
}

func parsePageMeta(metadata []byte, p *object.Page) error {
	if bytes.HasPrefix(metadata, []byte(pageMetaTOML)) {
		metadata = bytes.TrimPrefix(metadata, []byte(pageMetaTOML))
		metadata = bytes.TrimSpace(metadata)
		return toml.Unmarshal(metadata, p)
	}
	return ErrorPageMetaUnknownType
}

func fillOnePage(p *object.Page, info os.FileInfo) error {
	if p.Draft {
		return nil
	}
	p.ContentBytes = blackfriday.Run(p.RawBytes)
	var err error
	if p.CreateString == "" {
		p.CreateTime = info.ModTime()
	} else {
		if p.CreateTime, err = parseTimeString(p.CreateString); err != nil {
			return err
		}
	}
	if p.UpdateString == "" {
		p.UpdateTime = p.CreateTime
	} else {
		if p.UpdateTime, err = parseTimeString(p.UpdateString); err != nil {
			return err
		}
	}
	if p.CreateTime.Unix()-p.UpdateTime.Unix() > 0 {
		return ErrorTimeCreatedOverUpdated
	}
	p.Index = object.NewContentIndex(bytes.NewReader(p.ContentBytes))
	if p.Slug != "" {
		p.URL = p.Slug
		if !strings.HasSuffix(p.URL, ".html") {
			p.URL += ".html"
		}
	}
	return nil
}
