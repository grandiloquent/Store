package handlers

import (
	"net/http"
	"store/common"
)

func ServiceHandler(e *common.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		data := map[string]interface{}{
			"Title": "服务",
			"Debug": e.Debug}

		renderPage(w, "service.html", data, e.Debug)
	})
}
