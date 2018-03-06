package source

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/fuxiaohei/pugov1/module/object"
	"github.com/rs/zerolog/log"
)

// Copy run copied files to output dir
func Copy(s *object.Source, outputDir string) (int, int, error) {
	var (
		count int
		skip  int
	)
	for _, item := range s.CopyFiles {
		srcFile := item.File
		if item.SrcFile != "" {
			srcFile = item.SrcFile
		}
		dstFile := filepath.Join(outputDir, item.File)
		s.RenderedFiles = append(s.RenderedFiles, dstFile)

		info, _ := os.Stat(dstFile)
		if info != nil {
			if info.ModTime().Unix() == item.Info.ModTime().Unix() {
				log.Debug().Str("file", srcFile).Msg("copy-skip")
				skip++
				continue
			}
		}
		os.MkdirAll(filepath.Dir(dstFile), os.ModePerm)

		if err := copyFile(srcFile, dstFile, item.Info.ModTime()); err != nil {
			log.Warn().Str("file", srcFile).Err(err).Msg("copy-error")
			continue
		}
		log.Debug().Str("file", srcFile).Str("dest", dstFile).Msg("copy-ok")
		count++
	}
	return count, skip, nil
}

func copyFile(src string, dst string, modTime time.Time) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()
	d, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	if _, err := io.Copy(d, s); err != nil {
		d.Close()
		return err
	}
	os.Chtimes(dst, modTime, modTime)
	return d.Close()
}
