package admin

import (
	"net/http"

	"github.com/fuxiaohei/pugov1/module/config"
	"github.com/rs/zerolog/log"
)

func metaPageHandler(rw http.ResponseWriter, r *http.Request) {
	render(rw, "index.html", map[string]interface{}{
		"Meta": gConfig.Meta,
	})
}

func metaPostHandler(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	gConfig.Meta.Title = r.FormValue("title")
	gConfig.Meta.Subtitle = r.FormValue("subtitle")
	gConfig.Meta.Keyword = r.FormValue("keywords")
	gConfig.Meta.Root = r.FormValue("root")
	gConfig.Meta.Lang = r.FormValue("lang")
	gConfig.Meta.Desc = r.FormValue("desc")
	toFile := "test.toml"
	if err := config.Write(gConfig, toFile); err != nil {
		render(rw, "error.html", map[string]interface{}{
			"Error": err,
			"From":  "Saving Meta Error",
			"Back":  r.Referer(),
		})
		return
	}
	log.Info().Str("tofile", toFile).Msg("meta-write")
	rw.Header().Set("Location", "/_admin/meta?ok=1")
	rw.WriteHeader(302)
}
