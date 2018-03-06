package object

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// Post is a post
type Post struct {
	Title        string   `toml:"title" ini:"title"`
	Slug         string   `toml:"slug" ini:"slug"`
	Desc         string   `toml:"desc" ini:"desc"`
	CreateString string   `toml:"date" ini:"date"`
	UpdateString string   `toml:"update_date" ini:"update_date"`
	AuthorName   string   `toml:"author" ini:"author"`
	Thumb        string   `toml:"thumb" ini:"thumb"`
	Draft        bool     `toml:"draft" ini:"draft"`
	TagStrings   []string `toml:"tags" ini:"-"`

	MetaFormat    string `toml:"-" ini:"-"`
	SourceFile    string `toml:"-" ini:"-"`
	SourceRelpath string `toml:"-" ini:"-"`
	OutputFile    string `toml:"-" ini:"-"`
	URL           string `toml:"-" ini:"-"`

	Tags   []*PostTag      `toml:"-" ini:"-"`
	Author *Author         `toml:"-" ini:"-"`
	Index  []*ContentIndex `toml:"-" ini:"-"`

	CreateTime time.Time `toml:"-" ini:"-"`
	UpdateTime time.Time `toml:"-" ini:"-"`

	RawBytes     []byte `toml:"-" ini:"-"`
	ContentBytes []byte `toml:"-" ini:"-"`
	BriefBytes   []byte `toml:"-" ini:"-"`
}

// TagString return a string of post tags, joined with comma
func (p *Post) TagString() string {
	return strings.Join(p.TagStrings, ",")
}

// Created return post create time
func (p *Post) Created() time.Time {
	return p.CreateTime
}

// Updated return post updated time
func (p *Post) Updated() time.Time {
	return p.UpdateTime
}

// ContentHTML return content as HTML
func (p *Post) ContentHTML() template.HTML {
	return template.HTML(p.ContentBytes)
}

// BriefHTML return brief content as HTML
func (p *Post) BriefHTML() template.HTML {
	return template.HTML(p.BriefBytes)
}

// PostTag is tag in post
type PostTag struct {
	Name string
	URL  string
}

// PostList is list of posts with pager
type PostList struct {
	Pager *Pager
	Posts []*Post
}

// TagList is list of posts with Tag
type TagList struct {
	Tag      *PostTag
	Archives []*Archive
	Posts    []*Post
}

// ContentIndex is index of content
type ContentIndex struct {
	Level    int
	Title    string
	Archor   string
	Children []*ContentIndex
	Link     string
	Parent   *ContentIndex
}

// Print prints post indexs friendly
func (p *ContentIndex) Print() {
	fmt.Println(strings.Repeat("#", p.Level), p)
	for _, c := range p.Children {
		c.Print()
	}
}

// NewContentIndex parse reader byte to index
func NewContentIndex(r io.Reader) []*ContentIndex {
	var (
		z = html.NewTokenizer(r)

		currentLevel    int
		currentText     string
		currentLinkText string
		currentArchor   string
		nodeDeep        int

		indexs []*ContentIndex
	)
	for {
		token := z.Next()
		if token == html.ErrorToken {
			break
		}
		if token == html.EndTagToken {
			if nodeDeep == 1 && currentLevel > 0 {
				indexs = append(indexs, &ContentIndex{
					Level:  currentLevel,
					Title:  currentText,
					Link:   currentLinkText,
					Archor: currentArchor,
				})
				currentLevel = 0
				currentText = ""
				currentLinkText = ""
				currentArchor = ""
			}
			nodeDeep--
			continue
		}
		if token == html.StartTagToken {
			name, hasAttr := z.TagName()
			lv := parsePostIndexLevel(name)

			if lv > 0 {
				currentLevel = lv
				if hasAttr {
					for {
						k, v, isMore := z.TagAttr()
						if bytes.Equal(k, []byte("id")) {
							currentArchor = string(v)
						}
						if !isMore {
							break
						}
					}
				}
			}
			nodeDeep++

			if currentLevel > 0 && string(name) == "a" {
				if hasAttr {
					for {
						k, v, isMore := z.TagAttr()
						if bytes.Equal(k, []byte("href")) {
							currentLinkText = string(v)
						}
						if !isMore {
							break
						}
					}
				}
			}
		}
		if token == html.TextToken && currentLevel > 0 {
			currentText += string(z.Text())
		}
	}
	indexs = refreshIndexes(indexs)
	return indexs
}

func refreshIndexes(indexList []*ContentIndex) []*ContentIndex {
	var (
		list    []*ContentIndex
		lastIdx int
		lastN   *ContentIndex
	)
	for i, n := range indexList {
		if i == 0 {
			list = append(list, n)
			lastIdx = 0
			continue
		}
		lastN = list[lastIdx]
		if lastN.Level < n.Level {
			n.Parent = lastN
			lastN.Children = append(lastN.Children, n)
		} else {
			list = append(list, n)
			lastIdx++
		}
	}
	for _, n := range list {
		if len(n.Children) > 1 {
			n.Children = refreshIndexes(n.Children)
		}
	}
	return list
}

func parsePostIndexLevel(name []byte) int {
	if bytes.Equal(name, []byte("h1")) {
		return 1
	}
	if bytes.Equal(name, []byte("h2")) {
		return 2
	}
	if bytes.Equal(name, []byte("h3")) {
		return 3
	}
	if bytes.Equal(name, []byte("h4")) {
		return 4
	}
	if bytes.Equal(name, []byte("h5")) {
		return 5
	}
	if bytes.Equal(name, []byte("h6")) {
		return 6
	}
	return 0
}
