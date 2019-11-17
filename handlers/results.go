package handlers

import (
	"fmt"
	"net/http"
	"store/common"
)

func ResultsHandler(e *common.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		keyword := r.URL.Query().Get("keyword")
		fmt.Println(keyword)
	})
}
