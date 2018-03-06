package config

import (
	"errors"

	"github.com/go-ini/ini"

	"github.com/BurntSushi/toml"
	"github.com/fuxiaohei/pugov1/module/object"
)

type readHandler func(data []byte) (*object.Config, error)

var (
	// ErrorNoReadHandler means no handler to read config file
	ErrorNoReadHandler = errors.New("unsupported config file")
)

var (
	supportedHandlers = map[string]readHandler{
		".toml": readTOML,
		".ini":  readINI,
	}
)

func readTOML(data []byte) (*object.Config, error) {
	cfg := new(object.Config)
	return cfg, toml.Unmarshal(data, cfg)
}

func readINI(data []byte) (*object.Config, error) {
	iniObj, err := ini.Load(data)
	if err != nil {
		return nil, err
	}
	meta := new(object.Meta)
	if err = iniObj.Section("meta").MapTo(meta); err != nil {
		return nil, err
	}
	comments := new(object.Comment)
	if err = iniObj.Section("comment").MapTo(comments); err != nil {
		return nil, err
	}
	analytics := new(object.Analytics)
	if err = iniObj.Section("analytics").MapTo(analytics); err != nil {
		return nil, err
	}

	var navGroups object.NavGroup
	navKeys := iniObj.Section("nav").KeyStrings()
	for _, key := range navKeys {
		value := iniObj.Section("nav").Key(key).String()
		nav := new(object.Nav)
		if err = iniObj.Section("nav." + value).MapTo(nav); err != nil {
			return nil, err
		}
		if nav.Link == "" || nav.Title == "" {
			continue
		}
		navGroups = append(navGroups, nav)
	}

	var authGroups object.AuthorGroup
	authKeys := iniObj.Section("author").KeyStrings()
	for _, key := range authKeys {
		value := iniObj.Section("author").Key(key).String()
		author := new(object.Author)
		if err = iniObj.Section("author." + value).MapTo(author); err != nil {
			return nil, err
		}
		if author.Name == "" {
			continue
		}
		authGroups = append(authGroups, author)
	}

	return &object.Config{
		Meta:      meta,
		Comment:   comments,
		Analytics: analytics,
		Navs:      navGroups,
		Authors:   authGroups,
	}, nil
}
