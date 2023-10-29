package todo

import (
	"github.com/rs/zerolog/hlog"
	"gorm.io/gorm"
	"html/template"
	"net/http"

	"goji.io"
	"goji.io/pat"
)

type Todo struct {
	ID uint

	Title       string
	Description string
	Done        bool

	ProjectID uint
}

var templates *template.Template
var db *gorm.DB

func Handle(mux *goji.Mux, t *template.Template, d *gorm.DB) {
	templates = t
	db = d

	todoMux := goji.SubMux()

	todoMux.HandleFunc(pat.Get("/:id"), todoGetHandler)
	todoMux.HandleFunc(pat.Put("/:id"), todoPutHandler)
	todoMux.HandleFunc(pat.Delete("/:id"), todoDeleteHandler)

	mux.Handle(pat.New("/todo/*"), todoMux)
}

func todoDeleteHandler(w http.ResponseWriter, r *http.Request) {
	log := hlog.FromRequest(r)

	err := templates.ExecuteTemplate(w, ".html", nil)

	if err != nil {
		log.Err(err).Msg("template parse error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func todoPutHandler(w http.ResponseWriter, r *http.Request) {
	log := *hlog.FromRequest(r)

	err := r.ParseForm()
	if err != nil {
		log.Err(err).Msg("form parse error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = templates.ExecuteTemplate(w, ".html", nil)

	if err != nil {
		log.Err(err).Msg("template parse error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func todoGetHandler(w http.ResponseWriter, r *http.Request) {
	log := hlog.FromRequest(r)

	err := templates.ExecuteTemplate(w, ".html", nil)

	if err != nil {
		log.Err(err).Msg("template parse error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
