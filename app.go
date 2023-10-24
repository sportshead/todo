//go:generate ./tailwind.sh --minify --output ./static/tailwind.css
package main

import (
	"embed"
	"flag"
	"github.com/sportshead/todo/project"
	"html/template"
	"net/http"
	"os"
	"strings"
	"time"

	"goji.io"
	"goji.io/pat"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/CAFxX/httpcompression"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"

	"github.com/sportshead/todo/todo"
)

//go:embed templates/*.html
var templateFS embed.FS
var templates *template.Template

//go:embed static/tailwind.css
var staticFS embed.FS

func main() {
	log.Logger = log.
		Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).
		With().Caller().
		Logger()

	// https://stackoverflow.com/a/65215010
	templates = template.Must(template.New("").Funcs(template.FuncMap{
		"replaceNewlines": func(s string) template.HTML {
			return template.HTML(strings.Replace(template.HTMLEscapeString(s), "\n", "<br>", -1))
		},
	}).ParseFS(templateFS, "templates/*.html"))

	addr := flag.String("addr", "localhost:8080", "Address to listen on")
	file := flag.String("file", "data.db", "SQLite database file")
	flag.Parse()

	db, err := gorm.Open(sqlite.Open(*file), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Str("file", *file).Msg("failed to open database")
	}

	err = db.AutoMigrate(&todo.Todo{})
	if err != nil {
		log.Fatal().Err(err).Str("model", "todo").Msg("failed to migrate db")
		return
	}
	err = db.AutoMigrate(&project.Project{})
	if err != nil {
		log.Fatal().Err(err).Str("model", "project").Msg("failed to migrate db")
		return
	}
	log.Info().Str("file", *file).Msg("database migration complete")

	root := goji.NewMux()
	todo.Handle(root, templates)
	project.Handle(root, templates)

	root.Handle(pat.New("/"), http.RedirectHandler("/dashboard", http.StatusFound))
	root.Handle(pat.New("/static/*"), http.FileServer(http.FS(staticFS)))

	root.HandleFunc(pat.New("/dashboard"), dashboardHandler)

	root.Use(hlog.NewHandler(log.Logger))
	root.Use(hlog.MethodHandler("method"))
	root.Use(hlog.URLHandler("url"))
	root.Use(hlog.RefererHandler("referer"))

	compress, err := httpcompression.DefaultAdapter()
	if err != nil {
		log.Err(err).Msg("failed to create adapter for http compression")
	} else {
		root.Use(compress)
	}

	log.Info().Str("addr", *addr).Msg("server listening")
	err = http.ListenAndServe(*addr, root)
	log.Err(err).Msg("server stopped")
}
