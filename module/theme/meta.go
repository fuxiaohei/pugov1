package theme

import (
	"github.com/BurntSushi/toml"
	"github.com/fuxiaohei/pugov1/module/object"
)

// Meta is description of theme
type Meta struct {
	Name string   `toml:"name" ini:"name"`
	Repo string   `toml:"repo" ini:"repo"`
	URL  string   `toml:"url" ini:"url"`
	Date string   `toml:"date" ini:"date"`
	Desc string   `toml:"desc" ini:"desc"`
	Tags []string `toml:"tags" ini:"-"`

	MinVersion string `toml:"min_version" ini:"min_version"`

	Authors object.AuthorGroup `toml:"author" ini:"-"`
	Refs    []*metaRef         `toml:"ref" ini:"-"`

	License    string `toml:"license" ini:"license"`
	LicenseURL string `toml:"license_url" ini:"license_url"`
}

type metaRef struct {
	Name string `toml:"name" ini:"name"`
	URL  string `toml:"url" ini:"url"`
	Repo string `toml:"repo" ini:"repo"`
}

type readHandler func([]byte) (*Meta, error)

func readTOMLMeta(data []byte) (*Meta, error) {
	meta := new(Meta)
	return meta, toml.Unmarshal(data, meta)
}

var (
	metaReadHandlers = map[string]readHandler{
		".toml": readTOMLMeta,
	}
)
