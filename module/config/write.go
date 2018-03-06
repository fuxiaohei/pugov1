package config

import (
	"bytes"
	"errors"
	"io/ioutil"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/fuxiaohei/pugov1/module/object"
)

var (
	// ErrorWriteFileNotSet means not filename set when writing config
	ErrorWriteFileNotSet = errors.New("config-write-file-not-set")
	// ErrorWriteFileNotSupported means writing config to unknown file type
	ErrorWriteFileNotSupported = errors.New("config-write-file-not-supported")
)

// Write write config to file,
// If file is not set, write to cfg srcfile
func Write(cfg *object.Config, file string) error {
	toFile := file
	if toFile == "" {
		toFile = cfg.SrcFile
	}
	if toFile == "" {
		return ErrorWriteFileNotSet
	}
	suffix := filepath.Ext(toFile)
	switch suffix {
	case ".toml":
		return writeTOML(cfg, toFile)
	}
	return ErrorWriteFileNotSupported
}

func writeTOML(cfg *object.Config, file string) error {
	buf := bytes.NewBuffer(nil)
	encoder := toml.NewEncoder(buf)
	if err := encoder.Encode(cfg); err != nil {
		return err
	}
	return ioutil.WriteFile(file, buf.Bytes(), 0644)
}
