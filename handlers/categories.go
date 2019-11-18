package handlers

import (
	"net/http"
	"store/common"
	"text/template"
)

type Categories struct {
	Title string
	Debug bool
	Names []interface{}
}

func CategoriesHandler(e *common.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		names, err := e.DB.Fetch("select name from store_category")
		if err != nil {
			internalServerError(w, err)
			return
		}
		writeCategories(w, &Categories{
			Title: "品类",
			Debug: e.Debug,
			Names: names,
		})
	})
}
func writeCategories(w http.ResponseWriter, data *Categories) {
	t, err := template.ParseFiles("templates/categories.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		internalServerError(w, err)
		return
	}
	t.Execute(w, *data)
}
