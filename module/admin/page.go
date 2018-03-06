package admin

import (
	"errors"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/fuxiaohei/pugov1/module/config"
	"github.com/fuxiaohei/pugov1/module/object"
	"github.com/fuxiaohei/pugov1/module/source"
)

func pageListHandler(rw http.ResponseWriter, r *http.Request) {
	render(rw, "page.html", map[string]interface{}{
		"Pages": gSource.Pages,
	})
}

func pageEditHandler(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	slug := r.FormValue("slug-key")
	create := r.FormValue("created")
	ct, _ := strconv.ParseInt(create, 10, 64)
	var page *object.Page
	isNew := r.FormValue("new") != ""
	if isNew {
		page = &object.Page{
			Title:      "New Page",
			Slug:       "new-page-slug",
			CreateTime: time.Now(),
			Draft:      true,
			Template:   "page.html",
		}
	} else {
		page = source.QueryPage(gSource.Pages, slug, ct)
	}
	if page == nil {
		render(rw, "error.html", map[string]interface{}{
			"Error": errors.New("page is not found"),
			"From":  "Editing Page Error",
			"Back":  r.Referer(),
		})
		return
	}
	render(rw, "page-edit.html", map[string]interface{}{
		"Page":  page,
		"IsNew": isNew,
	})
}

func pageSaveHandler(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	isNew := r.FormValue("new") != ""
	var err error
	if isNew {
		var page *object.Page
		page, err = saveNewPage(r.Form)
		if err == nil {
			gSource.Pages = append([]*object.Page{page}, gSource.Pages...)
		}
	} else {
		slug := r.FormValue("slug-key")
		create := r.FormValue("created")
		ct, _ := strconv.ParseInt(create, 10, 64)
		page := source.QueryPage(gSource.Pages, slug, ct)
		if page == nil {
			render(rw, "error.html", map[string]interface{}{
				"Error": errors.New("page is not found"),
				"From":  "Saving Page Error",
				"Back":  r.Referer(),
			})
			return
		}
		err = saveExistPage(page, r.Form)
	}
	if err != nil {
		render(rw, "error.html", map[string]interface{}{
			"Error": err,
			"From":  "Saving Page Error",
			"Back":  r.Referer(),
		})
		return
	}
	rw.Header().Set("Location", "/_admin/pages")
	rw.WriteHeader(302)
}

func saveExistPage(page *object.Page, values url.Values) error {
	// create := values.Get("created")
	// ct, _ := strconv.ParseInt(create, 10, 64)
	title := values.Get("title")
	slug := values.Get("slug")
	if title == "" {
		return errors.New("title-or-slug-is-empty")
	}
	page.RawBytes = []byte(values.Get("content"))
	page.Desc = values.Get("desc")
	page.Title = title
	page.Slug = slug
	page.Draft = false
	page.UpdateTime = time.Now()
	page.UpdateString = page.UpdateTime.Format(source.TimeLayoutCommon)
	page.Template = values.Get("template")
	if values.Get("draft") != "" {
		page.Draft = true
	}
	// create correct url for page
	page.URL = strings.TrimSuffix(page.SourceRelpath, filepath.Ext(page.SourceRelpath)) + ".html"
	if page.Slug != "" {
		page.URL = page.Slug
		if !strings.HasSuffix(page.URL, ".html") {
			page.URL += ".html"
		}
	}
	return source.WritePage(page)
}

func saveNewPage(values url.Values) (*object.Page, error) {
	title := values.Get("title")
	slug := values.Get("slug")
	file := values.Get("file")
	if title == "" || file == "" {
		return nil, errors.New("title-or-or-file-is-empty")
	}
	p := &object.Page{
		Title:         title,
		Slug:          slug,
		SourceRelpath: file,
		SourceFile:    filepath.Join(gSource.PageDir, file),
		CreateTime:    time.Now(),
		RawBytes:      []byte(values.Get("content")),
		Template:      values.Get("template"),
		Desc:          values.Get("desc"),
	}
	p.CreateString = p.CreateTime.Format(source.TimeLayoutCommon)
	p.UpdateTime = p.CreateTime
	p.UpdateString = p.CreateString
	if owner := config.Owner(gConfig); owner != nil {
		p.AuthorName = owner.Name
	}
	// create correct url for page
	p.URL = strings.TrimSuffix(p.SourceRelpath, filepath.Ext(p.SourceRelpath)) + ".html"
	if p.Slug != "" {
		p.URL = p.Slug
		if !strings.HasSuffix(p.URL, ".html") {
			p.URL += ".html"
		}
	}
	return p, source.WritePage(p)
}
