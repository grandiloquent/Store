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
	FetchSearchKeywordsSQL  = "select * from store_fetch_keywords($1, $2,$3);"
)

func fetchSearch(e *common.Env, w http.ResponseWriter, r *http.Request) {
	limit := safeQueryInt(r, "limit", 6)
	offset := safeQueryInt(r, "offset", 0)
	sorttype := safeQueryInt(r, "sorttype", 2)

	s, err := fetchSearchKeywords(e, limit, offset, sorttype)

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
func fetchSearchKeywords(e *common.Env, limit, offset, sorttype int) ([]string, error) {
	rows, err := e.DB.Query(context.Background(), FetchSearchKeywordsSQL, limit, offset, sorttype)
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
		item := m.(map[string]interface{})
		search := item["search"].(string)
		raw := safeString(item["raw"])
		visits := item["visits"].(float64)
		popular := item["popular"].(float64)
		fmt.Println(item)
		batch.Queue(InsertSearchKeywordsSQL, search,
			raw,
			int64(visits),
			int64(popular),
		)
	}
	br := e.DB.SendBatch(context.Background(), batch)

	t, err := br.Exec()
	if err != nil {
		internalServerError(w, err)
		return
	}
	err = br.Close()
	if err != nil {
		internalServerError(w, err)
		return
	}
	writeCommandTag(t, w)
}
