package admin

import (
	"errors"
	"net/http"
	"strconv"
	"time"

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
	slug := r.FormValue("slug")
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
