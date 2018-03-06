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

func postPageHandler(rw http.ResponseWriter, r *http.Request) {
	render(rw, "post.html", map[string]interface{}{
		"Posts": gSource.Posts,
	})
}

func postEditHandler(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	slug := r.FormValue("slug")
	create := r.FormValue("created")
	ct, _ := strconv.ParseInt(create, 10, 64)
	var post *object.Post
	isNew := r.FormValue("new") != ""
	if isNew {
		post = &object.Post{
			Title:      "New Post",
			Slug:       "new-post-slug",
			CreateTime: time.Now(),
			Draft:      true,
		}
	} else {
		post = source.QueryPost(gSource.Posts, slug, ct)
	}
	if post == nil {
		render(rw, "error.html", map[string]interface{}{
			"Error": errors.New("post is not found"),
			"From":  "Editing Post Error",
			"Back":  r.Referer(),
		})
		return
	}
	render(rw, "post-edit.html", map[string]interface{}{
		"Post":  post,
		"IsNew": isNew,
	})
}

func postSaveHandler(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	isNew := r.FormValue("new") != ""
	var err error
	if isNew {
		err = saveNewPost(r.Form)
	} else {
		slug := r.FormValue("slug")
		create := r.FormValue("created")
		ct, _ := strconv.ParseInt(create, 10, 64)
		post := source.QueryPost(gSource.Posts, slug, ct)
		if post == nil {
			render(rw, "error.html", map[string]interface{}{
				"Error": errors.New("post is not found"),
				"From":  "Saving Post Error",
				"Back":  r.Referer(),
			})
			return
		}
		err = saveExistPost(post, r.Form)
	}
	if err != nil {
		render(rw, "error.html", map[string]interface{}{
			"Error": err,
			"From":  "Saving Post Error",
			"Back":  r.Referer(),
		})
		return
	}
	rw.Header().Set("Location", "/_admin/posts")
	rw.WriteHeader(302)
}

func saveExistPost(post *object.Post, values url.Values) error {
	create := values.Get("created")
	ct, _ := strconv.ParseInt(create, 10, 64)
	title := values.Get("title")
	slug := values.Get("slug")
	if title == "" || slug == "" {
		return errors.New("title-or-slug-is-empty")
	}
	if source.CheckPostSlugConflict(gSource.Posts, slug, ct) {
		return errors.New("post-slug-is-conflict")
	}
	post.RawBytes = []byte(values.Get("content"))
	post.Desc = values.Get("desc")
	post.Title = title
	post.TagStrings = strings.Split(values.Get("tags"), ",")
	post.Slug = slug
	post.Draft = false
	post.UpdateTime = time.Now()
	post.UpdateString = post.UpdateTime.Format(source.TimeLayoutCommon)
	if values.Get("draft") != "" {
		post.Draft = true
	}
	return source.WritePost(post)
}

func saveNewPost(values url.Values) error {
	title := values.Get("title")
	slug := values.Get("slug")
	file := values.Get("file")
	if title == "" || slug == "" || file == "" {
		return errors.New("title-or-slug-or-file-is-empty")
	}
	if source.CheckPostSlugConflict(gSource.Posts, slug, 0) {
		return errors.New("post-slug-is-conflict")
	}
	p := &object.Post{
		Title:         title,
		Slug:          slug,
		SourceRelpath: file,
		SourceFile:    filepath.Join(gSource.PostDir, file),
		CreateTime:    time.Now(),
		RawBytes:      []byte(values.Get("content")),
		TagStrings:    strings.Split(values.Get("tags"), ","),
		Desc:          values.Get("desc"),
	}
	p.CreateString = p.CreateTime.Format(source.TimeLayoutCommon)
	p.UpdateTime = p.CreateTime
	p.UpdateString = p.CreateString
	if owner := config.Owner(gConfig); owner != nil {
		p.AuthorName = owner.Name
	}
	return source.WritePost(p)
}
