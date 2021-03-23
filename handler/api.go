package handler

import (
	"net/http"

	"github.com/1garo/gopoc-executer/models"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

func executors(router chi.Router) {
	router.Get("/", executor)
	router.Post("/", executorReceive)
}

func executor(w http.ResponseWriter, r *http.Request) {
	//resOrdered.Users = append(resOrdered.Users, resAux.Users...)
	var rp models.Response
	rp.File = "aa"
	//   rp.= "retorno"
	if err := render.Render(w, r, rp); err != nil {
		render.Render(w, r, ErrorRenderer(err))
	}
}

func executorReceive(w http.ResponseWriter, r *http.Request) {
	rp := &models.Req{}
	if err := render.Bind(r, rp); err != nil {
		render.Render(w, r, ErrorRenderer(err))
	}

	if err := render.Render(w, r, rp); err != nil {
		render.Render(w, r, ErrorRenderer(err))
	}
}
