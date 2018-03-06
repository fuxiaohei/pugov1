package theme

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/fuxiaohei/pugov1/module/object"
	version "github.com/mcuadros/go-version"
)

var (
	reDefineTag   = regexp.MustCompile("{{ ?define \"([^\"]*)\" ?\"?([a-zA-Z0-9]*)?\"? ?}}")
	reTemplateTag = regexp.MustCompile("{{ ?template \"([^\"]*)\" ?([^ ]*)? ?}}")

	// ErrorMetaNotFound means theme's meta file is missing
	ErrorMetaNotFound = errors.New("need add theme meta file")
	// ErrorOutVersion means theme is out of version
	ErrorOutVersion = errors.New("theme need newer PuGo version")
)

type (
	// Theme object, maintains a sort of templates for whole site data
	Theme struct {
		Meta        *Meta
		pugoVersion string
		metaFile    string
		metaError   error
		dir         string
		lock        sync.Mutex
		funcMap     template.FuncMap
		templates   map[string]*template.Template
		extensions  []string

		cache               []*namedTemplate
		regularTemplateDefs []string
	}
	namedTemplate struct {
		Name string
		Src  string
	}
)

// Read read directory to load new theme
func Read(dir string, pugoVersion string) (*Theme, error) {
	theme := &Theme{
		dir:         dir,
		funcMap:     make(template.FuncMap),
		extensions:  []string{".html"},
		pugoVersion: pugoVersion,
	}
	theme.funcMap["HTML"] = func(v interface{}) template.HTML {
		if str, ok := v.(string); ok {
			return template.HTML(str)
		}
		if b, ok := v.([]byte); ok {
			return template.HTML(string(b))
		}
		return template.HTML(fmt.Sprintf("%v", v))
	}
	theme.funcMap["Include"] = func(values ...interface{}) template.HTML {
		var buf bytes.Buffer
		if len(values) < 2 {
			return template.HTML("<!-- include template without path or data -->")
		}
		var pathData []string
		for i, v := range values {
			if i < len(values)-1 {
				str, ok := v.(string)
				if !ok {
					return template.HTML("<!-- include template with non-string path -->")
				}
				pathData = append(pathData, str)
			}
		}
		tpl := path.Join(pathData...)
		if err := theme.Execute(&buf, tpl, values[len(values)-1]); err != nil {
			return template.HTML("<!-- template " + tpl + " error:" + err.Error() + "-->")
		}
		return template.HTML(string(buf.Bytes()))
	}
	theme.funcMap["String"] = func(value interface{}) string {
		switch value.(type) {
		case []byte:
			return string(value.([]byte))
		case template.HTML:
			return string(value.(template.HTML))
		default:
			return fmt.Sprint(value)
		}
	}
	return theme, theme.Load()
}

var (
	metaFiles = []string{
		"theme.toml",
	}
)

func (th *Theme) readMeta() error {
	for _, f := range metaFiles {
		file := path.Join(th.dir, f)
		if _, err := os.Stat(file); err != nil {
			continue
		}
		ext := filepath.Ext(file)
		handler := metaReadHandlers[ext]
		if handler == nil {
			continue
		}
		data, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}
		th.Meta, err = handler(data)
		if err != nil {
			return err
		}
		th.metaFile = file
		return nil
	}
	return ErrorMetaNotFound
}

// Func add template func to theme
func (th *Theme) Func(key string, fn interface{}) {
	th.funcMap[key] = fn
}

// Funcs return all template functions
func (th *Theme) Funcs() template.FuncMap {
	return th.funcMap
}

// Load loads templates
func (th *Theme) Load() error {
	if _, err := os.Stat(th.dir); err != nil {
		return err
	}
	if err := th.readMeta(); err != nil {
		return err
	}
	if err := th.checkMeta(); err != nil {
		return err
	}
	return th.loadTemplates()
}

func (th *Theme) checkMeta() error {
	if th.Meta.MinVersion != "" {
		if version.Compare(th.pugoVersion, th.Meta.MinVersion, "<") {
			return ErrorOutVersion
		}
	}
	return nil
}

// changes from https://github.com/go-macaron/renders/blob/master/renders.go#L43,
// thanks a lot
func (th *Theme) loadTemplates() error {
	th.lock.Lock()
	defer th.lock.Unlock()

	templates := make(map[string]*template.Template)

	err := filepath.Walk(th.dir, func(p string, fi os.FileInfo, err error) error {
		r, err := filepath.Rel(th.dir, p) // get relative path
		if err != nil {
			return err
		}
		ext := getExt(r)
		for _, extension := range th.extensions {
			if ext == extension {
				if err := th.add(p); err != nil {
					return err
				}
				for _, t := range th.regularTemplateDefs {
					found := false
					defineIdx := 0
					// From the beginning (which should) most specifc we look for definitions
					for _, nt := range th.cache {
						nt.Src = reDefineTag.ReplaceAllStringFunc(nt.Src, func(raw string) string {
							parsed := reDefineTag.FindStringSubmatch(raw)
							name := parsed[1]
							if name != t {
								return raw
							}
							// Don't touch the first definition
							if !found {
								found = true
								return raw
							}
							defineIdx++

							return fmt.Sprintf("{{ define \"%s_invalidated_#%d\" }}", name, defineIdx)
						})
					}
				}

				var (
					baseTmpl *template.Template
					i        int
				)

				for _, nt := range th.cache {
					var currentTmpl *template.Template
					if i == 0 {
						baseTmpl = template.New(nt.Name)
						currentTmpl = baseTmpl
					} else {
						currentTmpl = baseTmpl.New(nt.Name)
					}

					if _, err := currentTmpl.Funcs(th.funcMap).Parse(nt.Src); err != nil {
						return err
					}
					i++
				}
				tname := generateTemplateName(th.dir, p)
				templates[tname] = baseTmpl

				// Make sure we empty the cache between runs
				th.cache = th.cache[0:0]

				break
				//return nil
			}
		}
		return nil
	})
	th.templates = templates
	return err
}

func (th *Theme) add(path string) error {
	// Get file content
	tplSrc, err := getFileContent(path)
	if err != nil {
		return err
	}
	tplName := generateTemplateName(th.dir, path)
	// Make sure template is not already included
	alreadyIncluded := false
	for _, nt := range th.cache {
		if nt.Name == tplName {
			alreadyIncluded = true
			break
		}
	}
	if alreadyIncluded {
		return nil
	}

	// Add to the cache
	nt := &namedTemplate{
		Name: tplName,
		Src:  tplSrc,
	}
	th.cache = append(th.cache, nt)

	// Check for any template block
	for _, raw := range reTemplateTag.FindAllString(nt.Src, -1) {
		parsed := reTemplateTag.FindStringSubmatch(raw)
		templatePath := parsed[1]
		ext := getExt(templatePath)
		if !strings.Contains(templatePath, ext) {
			th.regularTemplateDefs = append(th.regularTemplateDefs, templatePath)
			continue
		}

		// Add this template and continue looking for more template blocks
		th.add(filepath.Join(th.dir, templatePath))
	}
	return nil
}

// Execute executes template by name with data,
// write into a Writer
func (th *Theme) Execute(w io.Writer, name string, data interface{}) error {
	tpl := th.Template(name)
	if tpl == nil {
		return fmt.Errorf("template '%s' is missing", name)
	}
	return tpl.ExecuteTemplate(w, name, data)
}

// StaticDir gets static dir in the theme
func (th *Theme) StaticDir() string {
	return path.Join(th.dir, th.Static())
}

// Dir get theme directory
func (th *Theme) Dir() string {
	return th.dir
}

// Static gets static dirname in the theme
func (th *Theme) Static() string {
	return "static"
}

// StaticFiles get static files in the theme
func (th *Theme) StaticFiles() ([]object.SourceItem, error) {
	var infos []object.SourceItem
	err := filepath.Walk(th.StaticDir(), func(fpath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		item := object.SourceItem{
			OpType:  object.SourceOpFile,
			Info:    info,
			File:    fpath,
			SrcFile: fpath,
		}
		item.File, _ = filepath.Rel(th.dir, fpath)
		infos = append(infos, item)
		return nil
	})
	return infos, err
}

// Template gets template by name
func (th *Theme) Template(name string) *template.Template {
	return th.templates[name]
}

func generateTemplateName(base, path string) string {
	//name := (r[0 : len(r)-len(ext)])
	return filepath.ToSlash(path[len(base)+1:])
}

func getFileContent(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	s := string(b)
	if len(s) < 1 {
		return "", errors.New("render: template file is empty")
	}
	return s, nil
}

func getExt(s string) string {
	if strings.Index(s, ".") == -1 {
		return ""
	}
	return "." + strings.Join(strings.Split(s, ".")[1:], ".")
}
