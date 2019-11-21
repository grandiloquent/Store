package handlers

import (
	"net/http"
	"store/common"
)

func CalculatorHandler(e *common.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		renderPage(w, "calculator.html", nil, e.Debug)
	})
}
