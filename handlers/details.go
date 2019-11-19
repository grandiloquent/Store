package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
	"store/common"
)

type Details struct {
	Title string
	Debug bool
}

func DetailsHandler(e *common.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		uid := vars["uid"]
		_ = uid
		renderPage(w, &Details{
			Debug: e.Debug,
		}, e.Debug)
	})
}
