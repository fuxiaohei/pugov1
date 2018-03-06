package packer

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

// Files pack files source code to gzip
func Files(files []string, dirs []string) map[string][]byte {
	for _, dir := range dirs {
		filelist, err := walkDir(dir)
		if err != nil {
			log.Warn().Str("dir", dir).Err(err).Msg("walk-dir-error")
			continue
		}
		files = append(files, filelist...)
	}
	res := make(map[string][]byte)
	for _, file := range files {
		fbytes, err := ioutil.ReadFile(file)
		if err != nil {
			log.Warn().Str("file", file).Err(err).Msg("read-file-error")
			continue
		}
		fbytes, err = Gzip(fbytes)
		if err != nil {
			log.Warn().Str("file", file).Err(err).Msg("gzip-file-error")
			continue
		}
		res[file] = fbytes
		log.Debug().Str("file", file).Msg("read-file")
	}
	log.Info().Int("files", len(res)).Msg("read-files")
	return res
}

func walkDir(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(fpath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		files = append(files, fpath)
		return nil
	})
	return files, err
}
