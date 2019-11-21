package handlers

import (
	"net/http"
	"store/common"
)

func CalculatorHandler(e *common.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"Debug": e.Debug,
			"Title": "计算器",
		}
		renderPage(w, "calculator.html", data, e.Debug)
	})
}
