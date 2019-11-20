package handlers

import (
	"fmt"
	"net/http"
	"store/common"
)

const (
	LikeSQL = "select * from store_list_like_results($1,$2,$3)"
)

func ResultsHandler(e *common.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		keyword := r.URL.Query().Get("keyword")
		like := "%" + keyword + "%"
		items, err := e.DB.Fetch(LikeSQL, like, 20, 0)
		if err != nil {
			internalServerError(w, err)
			return
		}
		m := map[string]interface{}{
			"Title": fmt.Sprintf("搜索 \"%s\"", keyword),
			"Debug": e.Debug,
			"Items": items,
		}
		renderPage(w, "results.html", m, e.Debug)

	})
}
