package source

import (
	"os"
	"path/filepath"

	"github.com/fuxiaohei/pugov1/module/i18n"
	"github.com/fuxiaohei/pugov1/module/object"
)

const (
	postExtension = ".md"
	pageExtension = ".md"
)

// Read read sources contents and maintains basic info in source object
func Read(postDir, pageDir, langDir string) (*object.Source, error) {
	sr := &object.Source{
		PostDir: postDir,
		PageDir: pageDir,
	}
	var err error
	sr.PostFiles, err = readPostFiles(postDir)
	if err != nil {
		return nil, err
	}
	sr.PageFiles, err = readPageFiles(pageDir)
	if err != nil {
		return nil, err
	}
	sr.I18ns, err = i18n.Read(langDir)
	if err != nil {
		return nil, err
	}
	return sr, nil
}

func readPostFiles(postDir string) ([]object.SourceItem, error) {
	var items []object.SourceItem
	err := filepath.Walk(postDir, func(fpath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		ext := filepath.Ext(fpath)
		item := object.SourceItem{
			SrcFile: fpath,
			File:    fpath,
			Info:    info,
		}
		if ext == postExtension {
			item.OpType = object.SourceOpPost
		} else {
			item.OpType = object.SourceOpFile
			relpath, _ := filepath.Rel(filepath.Dir(postDir), fpath)
			item.File = relpath
		}
		items = append(items, item)
		return nil
	})
	return items, err
}

func readPageFiles(pageDir string) ([]object.SourceItem, error) {
	var items []object.SourceItem
	err := filepath.Walk(pageDir, func(fpath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		ext := filepath.Ext(fpath)
		item := object.SourceItem{
			File: fpath,
			Info: info,
		}
		if ext == postExtension {
			item.OpType = object.SourceOpPage
		} else {
			item.OpType = object.SourceOpFile
		}
		items = append(items, item)
		return nil
	})
	return items, err
}
