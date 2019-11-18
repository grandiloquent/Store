package handlers

import (
	"net/http"
	"store/common"
	"text/template"
)

type Home struct {
	Title          string
	Debug          bool
	SearchHolder   string
	SearchKeywords []string
	Slide          []interface{}
}

func HomeHandler(e *common.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		searchKeywords, err := fetchSearchKeywords(e)
		if err != nil {
			internalServerError(w, err)
			return
		}
		slide, err := e.DB.Fetch("select filename from store_slide")
		if err != nil {
			internalServerError(w, err)
			return
		}
		writeHome(w, &Home{
			Title:          "淘货",
			Debug:          e.Debug,
			SearchHolder:   "精选好货",
			SearchKeywords: searchKeywords,
			Slide:          slide,
		})
	})
}
func writeHome(w http.ResponseWriter, data *Home) {
	t, err := template.ParseFiles("templates/home.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		internalServerError(w, err)
		return
	}
	t.Execute(w, *data)
}
