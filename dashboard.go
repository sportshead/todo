package main

import (
	"github.com/sportshead/todo/project"
	"github.com/sportshead/todo/todo"
	"net/http"

	"github.com/rs/zerolog/hlog"
)

type dashboardData struct {
	Projects []project.Project
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	log := hlog.FromRequest(r)

	err := templates.ExecuteTemplate(w, "dashboard.html", dashboardData{
		Projects: []project.Project{{
			ID:          1,
			Name:        "My Project",
			Description: "this is a description\nit can have multiple lines\nlike so",
			Archived:    false,
			Todos: []todo.Todo{
				{
					ID:          0,
					Title:       "finish this project",
					Description: "Very important!!!!!",
					Done:        false,
					ProjectID:   1,
				},
				{
					ID:          1,
					Title:       "execute order 66",
					Description: "using the youngling slayer 2000\nmultiple lines as well",
					Done:        false,
					ProjectID:   1,
				},
				{
					ID:          2,
					Title:       "become iron man",
					Description: "should be pretty easy",
					Done:        true,
					ProjectID:   1,
				},
				{
					ID:          3,
					Title:       "a super extremely long title name for no particular reason other then i feel like it",
					Description: "",
					Done:        false,
					ProjectID:   1,
				},
			},
		}},
	})

	if err != nil {
		log.Err(err).Msg("template parse error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
