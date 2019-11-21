package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx"
	"net/http"
	"store/common"
)

const (
	InsertSearchKeywordsSQL = "SELECT * FROM store_search_insert($1,$2,$3,$4)"
)

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
func fetchSearchKeywords(e *common.Env) ([]string, error) {
	rows, err := e.DB.Query(context.Background(), "select search from store_search limit 6")
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
func insertSearchKeywords(e *common.Env, w http.ResponseWriter, r *http.Request) {
	rItems, err := readData(e, r)
	if err != nil {
		internalServerError(w, err)
		return
	}
	items := (*rItems).([]interface{})

	batch := &pgx.Batch{}

	for _, m := range items {
		_ = m
		batch.Queue(InsertSearchKeywordsSQL)
	}
	results := e.DB.SendBatch(context.Background(), batch)
	t, err := results.Exec()
	if err != nil {
		internalServerError(w, err)
		return
	}
	c := t.RowsAffected()

	fmt.Println(items, c)

}
