package handlers

import (
	"net/http"
	"store/common"
)

func CartHandler(e *common.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{"Title": "购物车", "Debug": e.Debug}
		renderPage(w, "cart.html", data, e.Debug)
	})
}
