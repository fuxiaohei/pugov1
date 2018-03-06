package source

import (
	"os"
	"path/filepath"

	"github.com/fuxiaohei/pugov1/module/object"
	"github.com/rs/zerolog/log"
)

// Cleanup clean not rendered file
func Cleanup(s *object.Source, outputDir string) (int, error) {
	var count int
	files := make(map[string]bool)
	err := filepath.Walk(outputDir, func(fpath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		files[fpath] = false
		return nil
	})
	for _, file := range s.RenderedFiles {
		files[file] = true
	}
	for file, ok := range files {
		if !ok {
			os.Remove(file)
			log.Debug().Str("file", file).Msg("clean-file")
			count++
		}
	}
	return count, err
}
