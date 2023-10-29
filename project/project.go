package project

import (
	"github.com/rs/zerolog/hlog"
	"github.com/sportshead/todo/todo"
	"gorm.io/gorm"
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

type Data struct {
	ID uint

	Name        string
	Description string
	Archived    bool

	TotalTodos int
	DoneTodos  int
}

var templates *template.Template
var db *gorm.DB

func Handle(mux *goji.Mux, t *template.Template, d *gorm.DB) {
	templates = t
	db = d

	projectMux := goji.SubMux()

	projectMux.HandleFunc(pat.Get("/:id"), projectGetHandler)
	projectMux.HandleFunc(pat.Put("/:id"), projectPutHandler)
	projectMux.HandleFunc(pat.Delete("/:id"), projectDeleteHandler)

	projectMux.HandleFunc(pat.Get("/:id/todos"), projectGetTodosHandler)

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

// get project
func projectGetHandler(w http.ResponseWriter, r *http.Request) {
	log := hlog.FromRequest(r)

	data := Data{}
	result := db.Model(&Project{}).
		Joins("LEFT JOIN todos ON todos.project_id = projects.id").
		Select("projects.*, COUNT(todos.id) AS total_todos, SUM(todos.done) AS done_todos").
		Group("projects.id").
		Limit(1).
		Where("projects.id = ?", pat.Param(r, "id")).
		Find(&data)
	if result.Error != nil {
		log.Err(result.Error).Msg("db error")
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	if result.RowsAffected == 0 {
		log.Error().Msg("project not found")
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	err := templates.ExecuteTemplate(w, "project.html", data)

	if err != nil {
		log.Err(err).Msg("template parse error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func projectGetTodosHandler(w http.ResponseWriter, r *http.Request) {
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
