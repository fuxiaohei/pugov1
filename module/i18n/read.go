package i18n

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/rs/zerolog/log"
)

// Empty return empty i18n object
func Empty() *File {
	return &File{
		Lang:   "",
		File:   "",
		values: make(map[string]string),
	}
}

// Read read i18n files from a directory
func Read(dir string) (map[string]*File, error) {
	i18ns := make(map[string]*File)
	return i18ns, filepath.Walk(dir, func(fpath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		ext := filepath.Ext(fpath)
		if ext == ".toml" {
			fdata, _ := ioutil.ReadFile(fpath)
			f, err := readTOML(fdata, fpath)
			if err != nil {
				log.Warn().Err(err).Str("file", fpath).Msg("read-lang-error")
			}
			if old := i18ns[f.Lang]; old != nil {
				log.Warn().Err(err).Str("lang", f.Lang).Strs("file", []string{
					f.File, old.File,
				}).Msg("read-lang-conflict")
				return nil
			}
			i18ns[f.Lang] = f
		}
		return nil
	})
}

// File is file data from i18n file
type File struct {
	Lang   string
	File   string
	values map[string]string
}

// Tr transform i18n key to proper value
func (f *File) Tr(key string, defaults ...string) string {
	layout := f.values[key]
	if layout == "" {
		if len(defaults) > 0 {
			return defaults[0]
		}
	}
	return ""
}

func readTOML(data []byte, fromFile string) (*File, error) {
	m := make(map[string]interface{})
	if err := toml.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	f := &File{
		Lang:   fmt.Sprint(m["lang"]),
		File:   fromFile,
		values: make(map[string]string),
	}
	for k, v := range m {
		if mdata, ok := v.(map[string]interface{}); ok {
			for k2, v2 := range mdata {
				f.values[k+"."+k2] = fmt.Sprint(v2)
			}
		}
	}
	return f, nil
}
