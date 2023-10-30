package todo

import (
	"errors"
	"github.com/rs/zerolog/hlog"
	"gorm.io/gorm"
	"html/template"
	"net/http"
	"strconv"

	"goji.io"
	"goji.io/pat"
)

type Todo struct {
	ID uint

	Title       string
	Description string
	Done        bool

	ProjectID uint

	Deleted gorm.DeletedAt
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

	mux.HandleFunc(pat.Post("/todo"), rootPostHandler)
	mux.Handle(pat.New("/todo/*"), todoMux)
}

func todoDeleteHandler(w http.ResponseWriter, r *http.Request) {
	log := hlog.FromRequest(r)

	id, err := strconv.Atoi(pat.Param(r, "id"))
	if err != nil {
		log.Err(err).Str("id", pat.Param(r, "id")).Msg("invalid id")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := db.Delete(&Todo{}, id)
	if result.Error != nil {
		log.Err(result.Error).Msg("db error")
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func todoPutHandler(w http.ResponseWriter, r *http.Request) {
	log := *hlog.FromRequest(r)

	err := r.ParseForm()
	if err != nil {
		log.Err(err).Msg("form parse error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := strconv.Atoi(pat.Param(r, "id"))
	if err != nil {
		log.Err(err).Str("id", pat.Param(r, "id")).Msg("invalid id")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	projectId, err := strconv.Atoi(r.FormValue("projectId"))
	if err != nil {
		log.Err(err).Str("id", r.FormValue("projectId")).Msg("invalid project id")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data := Todo{
		ID:          uint(id),
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		Done:        r.FormValue("done") == "on",
		ProjectID:   uint(projectId),
	}
	result := db.Save(data)
	if result.Error != nil {
		log.Err(result.Error).Msg("db error")
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	err = templates.ExecuteTemplate(w, "todo.html", data)

	if err != nil {
		log.Err(err).Msg("template parse error")
		return
	}
}

func todoGetHandler(w http.ResponseWriter, r *http.Request) {
	log := hlog.FromRequest(r)

	id, err := strconv.Atoi(pat.Param(r, "id"))
	if err != nil {
		log.Err(err).Str("id", pat.Param(r, "id")).Msg("invalid id")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data := Todo{ID: uint(id)}
	result := db.First(&data)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			http.Error(w, result.Error.Error(), http.StatusNotFound)
			return
		}
		log.Err(result.Error).Msg("db error")
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	err = templates.ExecuteTemplate(w, "todo.html", data)

	if err != nil {
		log.Err(err).Msg("template parse error")
		return
	}
}

func rootPostHandler(w http.ResponseWriter, r *http.Request) {
	log := hlog.FromRequest(r)

	err := r.ParseForm()
	if err != nil {
		log.Err(err).Msg("form parse error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	projectId, err := strconv.Atoi(r.FormValue("projectId"))
	if err != nil {
		log.Err(err).Str("id", r.FormValue("projectId")).Msg("invalid project id")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data := Todo{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		Done:        r.FormValue("done") == "on",
		ProjectID:   uint(projectId),
	}
	result := db.Create(&data)
	if result.Error != nil {
		log.Err(result.Error).Msg("db error")
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = templates.ExecuteTemplate(w, "todo.html", data)

	if err != nil {
		log.Err(err).Msg("template parse error")
		return
	}
}
