package admin

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/fuxiaohei/pugov1/module/object"
	"github.com/fuxiaohei/pugov1/module/theme"
	"github.com/rs/zerolog/log"
)

var (
	the     *theme.Theme
	gSource *object.Source
	gConfig *object.Config
)

// Init init admin panels
func Init(themeDir string, cfg *object.Config, s *object.Source) {
	log.Info().Msg("admin-init")
	t, err := theme.Read(themeDir, "")
	if err != nil {
		log.Fatal().Err(err).Msg("admin-theme-error")
	}
	the = t
	log.Debug().Msg("admin-theme-ok")

	gSource = s
	gConfig = cfg
}

// Handler return http handler function for admin panels
func Handler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if the == nil {
			rw.WriteHeader(500)
			rw.Write([]byte("theme-error"))
			return
		}
		if r.Method == "POST" {
			postHandlers(rw, r)
			return
		}
		if r.Method != "GET" {
			rw.WriteHeader(405)
			return
		}
		if findFile(rw, r) {
			return
		}
		path := strings.TrimPrefix(r.URL.Path, "/_admin/")
		switch path {
		case "meta":
			metaPageHandler(rw, r)
		case "posts":
			postPageHandler(rw, r)
		case "posts/edit":
			postEditHandler(rw, r)
		case "pages":
			pageListHandler(rw, r)
		case "pages/edit":
			pageEditHandler(rw, r)
		default:
			rw.Header().Set("Location", "/_admin/meta")
			rw.WriteHeader(302)
		}
	}
}

func postHandlers(rw http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/_admin/")
	switch path {
	case "meta":
		metaPostHandler(rw, r)
	case "posts/edit":
		postSaveHandler(rw, r)
	default:
		rw.Header().Set("Location", "/_admin/meta")
		rw.WriteHeader(302)
	}
}

func findFile(rw http.ResponseWriter, r *http.Request) bool {
	relpath, _ := filepath.Rel("/_admin", r.RequestURI)
	file := filepath.Join(the.Dir(), relpath)
	if info, err := os.Stat(file); err == nil && !info.IsDir() {
		http.ServeFile(rw, r, file)
		return true
	}
	return false
}

func render(rw http.ResponseWriter, tpl string, values interface{}) {
	if err := the.Execute(rw, tpl, values); err != nil {
		log.Error().Err(err).Str("tpl", tpl).Msg("render-error")
		rw.WriteHeader(500)
		rw.Write([]byte("theme-error"))
	}
}
