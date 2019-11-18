package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"store/common"
)

func DetailsHandler(e *common.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		uid := vars["uid"]
		fmt.Println(uid)
	})
}
