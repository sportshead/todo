package project

import (
	"github.com/rs/zerolog/hlog"
	"github.com/sportshead/todo/todo"
	"html/template"
	"net/http"

	"goji.io"
	"goji.io/pat"
)

type Project struct {
	ID uint

	Name        string
	Description string
	Archived    bool

	Todos []todo.Todo
}

var templates *template.Template

func Handle(mux *goji.Mux, t *template.Template) {
	templates = t

	projectMux := goji.SubMux()

	projectMux.HandleFunc(pat.Get("/:id"), projectGetHandler)
	projectMux.HandleFunc(pat.Put("/:id"), projectPutHandler)
	projectMux.HandleFunc(pat.Delete("/:id"), projectDeleteHandler)

	mux.HandleFunc(pat.Get("/project"), rootGetHandler)
	mux.HandleFunc(pat.Post("/project"), rootPostHandler)
	mux.Handle(pat.New("/project/*"), projectMux)
}

func projectDeleteHandler(w http.ResponseWriter, r *http.Request) {
	log := hlog.FromRequest(r)

	err := templates.ExecuteTemplate(w, ".html", nil)

	if err != nil {
		log.Err(err).Msg("template parse error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func projectPutHandler(w http.ResponseWriter, r *http.Request) {
	log := hlog.FromRequest(r)

	err := templates.ExecuteTemplate(w, ".html", nil)

	if err != nil {
		log.Err(err).Msg("template parse error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// get project and todos
func projectGetHandler(w http.ResponseWriter, r *http.Request) {
	log := hlog.FromRequest(r)

	err := templates.ExecuteTemplate(w, ".html", nil)

	if err != nil {
		log.Err(err).Msg("template parse error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// new project
func rootPostHandler(w http.ResponseWriter, r *http.Request) {
	log := hlog.FromRequest(r)

	err := templates.ExecuteTemplate(w, ".html", nil)

	if err != nil {
		log.Err(err).Msg("template parse error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// list projects
func rootGetHandler(w http.ResponseWriter, r *http.Request) {
	log := hlog.FromRequest(r)

	err := templates.ExecuteTemplate(w, ".html", nil)

	if err != nil {
		log.Err(err).Msg("template parse error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
