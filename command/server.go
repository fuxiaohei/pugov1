package command

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/fuxiaohei/pugov1/module/admin"
	"github.com/fuxiaohei/pugov1/module/config"
	"github.com/fuxiaohei/pugov1/module/source"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli"
)

// Server is command of 'server'
var Server = cli.Command{
	Name:  "server",
	Usage: "run http server on output directory",
	Flags: append(commonFlags, cli.BoolFlag{
		Name:  "admin",
		Usage: "open admin panel pages",
	}, cli.StringFlag{
		Name:  "port",
		Usage: "set http server listening port",
		Value: defaultPort,
	}),
	Action: func(ctx *cli.Context) error {
		setLogLevel(ctx)
		return serverFunc(ctx)
	},
}

var defaultPort = "0.0.0.0:9899"
var fs http.Handler

func serverFunc(ctx *cli.Context) error {
	log.Info().Str("directory", outputDir).Str("port", defaultPort).Msg("http-listen")

	var adminFn http.Handler
	if ctx.Bool("admin") {
		cfg, err := config.Read()
		if err != nil {
			return err
		}
		s, err := source.Read(
			filepath.Join(sourceDir, postDir),
			filepath.Join(sourceDir, pageDir),
			filepath.Join(sourceDir, langDir),
		)
		if err != nil {
			return err
		}
		if err = source.Parse(s, true); err != nil {
			return err
		}
		log.Debug().Msg("source-parsed")
		admin.Init(filepath.Join(sourceDir, adminDir), cfg, s)
		adminFn = admin.Handler()
		go watchOnce([]string{sourceDir, themeDir}, ctx)
	}

	fs = http.FileServer(http.Dir(outputDir))
	var fn = func(rw http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.RequestURI, "/_admin") && adminFn != nil {
			adminFn.ServeHTTP(rw, r)
			return
		}
		fs.ServeHTTP(rw, r)
	}

	http.Handle("/", buildLogHandler(http.HandlerFunc(fn)))
	return http.ListenAndServe(ctx.String("port"), nil)
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.ResponseWriter.WriteHeader(code)
	rw.status = code
}

func buildLogHandler(fs http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		w := &responseWriter{ResponseWriter: rw}
		fs.ServeHTTP(w, r)
		if w.status >= 400 {
			log.Warn().Int("status", w.status).Str("method", r.Method).Str("url", r.RequestURI).Msg("req-fail")
			return
		}
		log.Debug().Int("status", w.status).Str("method", r.Method).Str("url", r.RequestURI).Msg("req-ok")
	})
}
