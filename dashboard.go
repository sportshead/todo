package main

import (
	"github.com/sportshead/todo/project"
	"gorm.io/gorm"
	"net/http"

	"github.com/rs/zerolog/hlog"
)

type dashboardData struct {
	ShowArchived bool
	Projects     *[]project.Data

	TotalTodos int
	DoneTodos  int
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	log := hlog.FromRequest(r)

	err := r.ParseForm()
	if err != nil {
		log.Err(err).Msg("error parsing form")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := dashboardData{
		ShowArchived: r.Form.Has("showArchived"),
		Projects:     new([]project.Data),
	}

	var result *gorm.DB
	if data.ShowArchived {
		result = db.Model(&project.Project{}).
			Joins("LEFT JOIN todos ON todos.project_id = projects.id").
			Select("projects.*, COUNT(todos.id) AS total_todos, SUM(todos.done) AS done_todos").
			Group("projects.id").
			Find(data.Projects)
	} else {
		result = db.Model(&project.Project{}).
			Joins("LEFT JOIN todos ON todos.project_id = projects.id").
			Select("projects.*, COUNT(todos.id) AS total_todos, SUM(todos.done) AS done_todos").
			Where("projects.archived = 0").
			Group("projects.id").
			Find(data.Projects)
	}
	if result.Error != nil {
		log.Err(result.Error).Msg("db error")
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	for _, proj := range *data.Projects {
		data.TotalTodos += proj.TotalTodos
		data.DoneTodos += proj.DoneTodos
	}

	log.Debug().Int64("rows", result.RowsAffected).Int("total_todos", data.TotalTodos).Int("done_todos", data.DoneTodos).Msg("got dashboard data")

	err = templates.ExecuteTemplate(w, "dashboard.html", data)

	if err != nil {
		log.Err(err).Msg("template parse error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
