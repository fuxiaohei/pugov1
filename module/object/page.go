package object

import (
	"html/template"
	"time"
)

// Page is a page
type Page struct {
	Title        string                 `toml:"title" ini:"title"`
	Slug         string                 `toml:"slug" ini:"slug"`
	Desc         string                 `toml:"desc" ini:"desc"`
	CreateString string                 `toml:"date" ini:"date"`
	UpdateString string                 `toml:"update_date" ini:"update_date"`
	AuthorName   string                 `toml:"author" ini:"author"`
	NavHover     string                 `toml:"hover" ini:"hover"`
	Template     string                 `toml:"template" ini:"template"`
	Lang         string                 `toml:"lang" ini:"lang"`
	Meta         map[string]interface{} `toml:"meta" ini:"-"`
	Sort         int                    `toml:"sort" ini:"sort"`
	Author       *Author                `toml:"-" ini:"-"`
	Draft        bool                   `toml:"draft" ini:"draft"`
	Index        []*ContentIndex        `toml:"-" ini:"-"`

	MetaFormat    string `toml:"-" ini:"-"`
	SourceFile    string `toml:"-" ini:"-"`
	SourceRelpath string `toml:"-" ini:"-"`
	OutputFile    string `toml:"-" ini:"-"`
	URL           string `toml:"-" ini:"-"`

	CreateTime time.Time `toml:"-" ini:"-"`
	UpdateTime time.Time `toml:"-" ini:"-"`

	RawBytes     []byte `toml:"-" ini:"-"`
	ContentBytes []byte `toml:"-" ini:"-"`
}

// Created return page created time
func (p *Page) Created() time.Time {
	return p.CreateTime
}

// Updated return page updated time
func (p *Page) Updated() time.Time {
	return p.UpdateTime
}

// ContentHTML return page content as HTML type
func (p *Page) ContentHTML() template.HTML {
	return template.HTML(p.ContentBytes)
}
