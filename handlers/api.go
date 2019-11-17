package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"store/common"
)

func ApiSearchHandler(e *common.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method := r.URL.Query().Get("method")
		if r.Method == "POST" {
			switch method {
			case "insert":
				insertSearch(e, w, r)
				return
			}
		} else {
			switch method {
			case "fetch":
				fetchSearch(e, w, r)
				return
			}
		}

		notFound(w)
	})
}
func fetchSearch(e *common.Env, w http.ResponseWriter, r *http.Request) {
	s, err := fetchSearchKeywords(e)
	if err != nil {
		internalServerError(w, err)
		return
	}
	obj, err := json.Marshal(s)
	if err != nil {
		internalServerError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(obj)
}
func insertSearch(e *common.Env, w http.ResponseWriter, r *http.Request) {
	items, err := readData(e, w, r)
	if err != nil {
		internalServerError(w, err)
		return
	}

	s := joinArray(*items)

	t, err := e.DB.Exec("select * from store_search_insert($1)", s)
	if err != nil {
		internalServerError(w, err)
		return
	}
	w.Write([]byte(fmt.Sprintf("%d", t.RowsAffected())))
}
func fetchSearchKeywords(e *common.Env) ([]string, error) {
	rows, err := e.DB.Query("select search from store_search limit 6")
	if err != nil {
		return nil, err
	}
	var items []string

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, err
		}
		items = append(items, values[0].(string))
	}
	return items, nil
}
