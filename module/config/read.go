package config

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/fuxiaohei/pugov1/module/object"
)

// Read read config data from file
func Read() (*object.Config, error) {
	file, ext, err := detectFile(supportedExtensions)
	if err != nil {
		return nil, err
	}
	handler := supportedHandlers[ext]
	if handler == nil {
		return nil, ErrorNoReadHandler
	}
	fileData, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	cfg, err := handler(fileData)
	if err != nil {
		return nil, err
	}
	cfg.SrcFile = file
	return cfg, nil
}

var (
	supportedExtensions = []string{".toml", ".ini"}
)

var (
	// ErrorNoConfigFile means config file is not found
	ErrorNoConfigFile = errors.New("config file is not found")
)

func detectFile(exts []string) (string, string, error) {
	filename := "config"
	for _, ext := range exts {
		file := filename + ext
		if _, err := os.Stat(file); err != nil {
			continue
		}
		return file, ext, nil
	}
	return "", "", ErrorNoConfigFile
}

// Owner get first author in config
func Owner(cfg *object.Config) *object.Author {
	if len(cfg.Authors) > 0 {
		return cfg.Authors[0]
	}
	return nil
}
