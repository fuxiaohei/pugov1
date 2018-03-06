package source

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/fuxiaohei/pugov1/module/object"
	"github.com/go-ini/ini"
	"github.com/rs/zerolog/log"
	blackfriday "gopkg.in/russross/blackfriday.v2"
)

func parsePosts(s *object.Source, withDraft bool) ([]*object.Post, error) {
	var posts []*object.Post
	for _, item := range s.PostFiles {
		if item.OpType == object.SourceOpFile {
			s.CopyFiles = append(s.CopyFiles, item)
			continue
		}
		post, err := parseOnePost(item.File, item.Info)
		if err != nil {
			log.Warn().Str("src", item.File).Err(err).Msg("read-post-error")
			continue
		}
		post.SourceRelpath, _ = filepath.Rel(s.PostDir, post.SourceFile)
		if post.Draft {
			log.Warn().Str("src", item.File).Msg("post-draft")
			if !withDraft {
				continue
			}
		}
		posts = append(posts, post)
	}
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].CreateTime.Unix() > posts[j].CreateTime.Unix()
	})
	return posts, nil
}

const (
	postSeparator      = "```"
	postMetaTOML       = "toml"
	postMetaINI        = "ini"
	postBriefSeparator = "<!--more-->"
)

var (
	// ErrorPostSectionsInvalid means post cant not seperate to current sections for parsing
	ErrorPostSectionsInvalid = errors.New("post sections are invalid")
	// ErrorPostMetaUnknownType means post metadata cant be unmarshaled with proper type
	ErrorPostMetaUnknownType = errors.New("post metadata is unknown type")
)

func parseOnePost(file string, fileInfo os.FileInfo) (*object.Post, error) {
	fData, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	sections := bytes.SplitN(fData, []byte(postSeparator), 3)
	if len(sections) < 3 {
		return nil, ErrorPostSectionsInvalid
	}
	p := new(object.Post)
	if err := parsePostMeta(sections[1], p); err != nil {
		return nil, err
	}
	p.RawBytes = bytes.TrimSpace(sections[2])
	p.SourceFile = file
	return p, fillOnePost(p, fileInfo, []byte(postBriefSeparator))
}

func parsePostMeta(metadata []byte, p *object.Post) error {
	if bytes.HasPrefix(metadata, []byte(postMetaTOML)) {
		p.MetaFormat = "toml"
		metadata = bytes.TrimPrefix(metadata, []byte(postMetaTOML))
		metadata = bytes.TrimSpace(metadata)
		return toml.Unmarshal(metadata, p)
	}
	if bytes.HasPrefix(metadata, []byte(postMetaINI)) {
		p.MetaFormat = "ini"
		metadata = bytes.TrimPrefix(metadata, []byte(postMetaINI))
		metadata = bytes.TrimSpace(metadata)
		iniObj, err := ini.Load(metadata)
		if err != nil {
			return err
		}
		if err = iniObj.MapTo(p); err != nil {
			return err
		}
		tagStr := iniObj.Section("").Key("tags").String()
		p.TagStrings = strings.Split(tagStr, ",")
		return nil
	}
	return ErrorPostMetaUnknownType
}

func fillOnePost(p *object.Post, info os.FileInfo, briefSeparator []byte) error {
	p.ContentBytes = blackfriday.Run(p.RawBytes)
	if len(briefSeparator) > 0 {
		sections := bytes.Split(p.RawBytes, briefSeparator)
		if len(sections) > 1 {
			p.BriefBytes = blackfriday.Run(sections[0])
		}
	}
	if len(p.BriefBytes) == 0 {
		p.BriefBytes = p.ContentBytes
	}
	for _, name := range p.TagStrings {
		p.Tags = append(p.Tags, &object.PostTag{
			Name: name,
			URL:  fmt.Sprintf("/tags/%s.html", name),
		})
	}
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
	p.URL = fmt.Sprintf("/%d/%d/%d/%s.html", p.CreateTime.Year(), p.CreateTime.Month(), p.CreateTime.Day(), p.Slug)
	return nil
}

var (
	// ErrorTimeParseEmpty means parsing empty string
	ErrorTimeParseEmpty = errors.New("parse empty time string")
	// ErrorTimeUnknownLayout means time string layout is not supported
	ErrorTimeUnknownLayout = errors.New("parse time but unknown layout")
	// ErrorTimeCreatedOverUpdated means create time is newer than updated time
	ErrorTimeCreatedOverUpdated = errors.New("created time is over updated time")
)

const (
	// TimeLayoutShort use date to print time
	TimeLayoutShort = "2006-01-02"
	// TimeLayoutCommon use date and daytime to print time
	TimeLayoutCommon = "2006-01-02 15:04"
	// TimeLayoutLong use date and daytime with seconds to print time
	TimeLayoutLong = "2006-01-02 15:04:05"
)

func parseTimeString(timeStr string) (time.Time, error) {
	timeStr = strings.TrimSpace(timeStr)
	if len(timeStr) == 0 {
		return time.Time{}, ErrorTimeParseEmpty
	}
	if len(timeStr) == len("2006-01-02") {
		return time.Parse("2006-01-02", timeStr)
	}
	if len(timeStr) == len("2006-01-02 15:04") {
		return time.Parse("2006-01-02 15:04", timeStr)
	}
	if len(timeStr) == len("2006-01-02 15:04:05") {
		return time.Parse("2006-01-02 15:04:05", timeStr)
	}
	return time.Time{}, ErrorTimeUnknownLayout
}

// WritePost write post to p.SourceFile
func WritePost(p *object.Post) error {
	buf := bytes.NewBuffer(nil)
	if p.MetaFormat == "toml" || p.MetaFormat == "" {
		buf.WriteString("```toml\n")
		encoder := toml.NewEncoder(buf)
		if err := encoder.Encode(p); err != nil {
			return err
		}
		buf.WriteString("```")
	}
	if p.MetaFormat == "ini" {
		buf.WriteString("```ini\n")
		iniObj := ini.Empty()
		if err := ini.ReflectFrom(iniObj, p); err != nil {
			return err
		}
		iniObj.Section("").Key("tags").SetValue(p.TagString())
		if _, err := iniObj.WriteToIndent(buf, "  "); err != nil {
			return err
		}
		buf.WriteString("```")
	}
	buf.WriteString("\n\n")
	buf.Write(p.RawBytes)
	toFile := p.SourceFile
	os.MkdirAll(filepath.Dir(toFile), 0755)
	return ioutil.WriteFile(toFile, buf.Bytes(), 0644)
}
