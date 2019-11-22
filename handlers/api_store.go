package handlers

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx"
	"github.com/pkg/errors"
	"net/http"
	"store/common"
)

func fetchStoreList(e *common.Env, w http.ResponseWriter, r *http.Request) {
	limit := safeQueryInt(r, "limit", 10)
	offset := safeQueryInt(r, "offset", 10)
	items, err := e.DB.Fetch(ListStoreSQL, limit, offset)
	if err != nil {
		internalServerError(w, err)
		return
	}
	var rows []interface{}
	for _, i := range items {
		item := i.([]interface{})
		item[2] = stringPrice(item[2])
		rows = append(rows, item)
	}
	buf, err := json.Marshal(rows)
	if err != nil {
		internalServerError(w, err)
		return
	}
	w.Write(buf)
}

func insertStore(e *common.Env, w http.ResponseWriter, r *http.Request) {
	items, err := readData(e, r)
	if err != nil {
		internalServerError(w, err)
		return
	}
	rows, ok := (*items).(map[string]interface{})
	if !ok {
		badRequest(w)
		return
	}
	// -----------------------------------
	title := rows["title"]
	if isWhiteSpaceString(title) {
		badRequest(w)
		return
	}
	price, err := toFloat(rows["price"])
	if err != nil {
		internalServerError(w, err)
		return
	}
	thumbnail := rows["thumbnail"]
	details := rows["details"]
	specification := rows["specification"]
	service := rows["service"]
	properties := joinArray(rows["properties"])
	showcases := joinArray(rows["showcases"])
	// -----------------------------------
	var uid string
	err = e.DB.QueryRow(context.Background(), InsertStoreSQL,
		title,
		price,
		thumbnail,
		details,
		specification,
		service,
		properties,
		showcases).Scan(&uid);
	if err != nil {
		internalServerError(w, err)
		return
	}
	// -----------------------------------
	w.Write([]byte(uid))
}
func insertStores(e *common.Env, w http.ResponseWriter, r *http.Request) {
	err := ensureConnection(e)
	if err != nil {
		internalServerError(w, err)
		return
	}
	// -----------------------------------

	rItems, err := readData(e, r)
	if err != nil {
		internalServerError(w, err)
		return
	}
	items, ok := (*rItems).([]interface{})
	if !ok {
		internalServerError(w, errors.New("invalid data"))
		return
	}

	batch := &pgx.Batch{}

	for _, v := range items {
		m, ok := v.(map[string]interface{})

		var uidValue interface{}
		uid, ok := m["uid"].(string)
		if !ok {
			uidValue = nil
		} else {
			uidValue = uid
		}
		title := m["title"]
		if isWhiteSpaceString(title) {
			badRequest(w)
			return
		}
		price, err := toFloat(m["price"])
		if err != nil {
			internalServerError(w, err)
			return
		}
		thumbnail := m["thumbnail"]
		details := m["details"]
		specification := m["specification"]
		service := m["service"]
		properties := joinArray(m["properties"])
		showcases := joinArray(m["showcases"])

		batch.Queue(InsertStoreSQL,
			uidValue,
			title,
			price,
			thumbnail,
			details,
			specification,
			service,
			properties,
			showcases,
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
	//items, err := readData(e, r)
	//if err != nil {
	//	internalServerError(w, err)
	//	return
	//}
	//rows, ok := (*items).(map[string]interface{})
	//if !ok {
	//	badRequest(w)
	//	return
	//}
	//// -----------------------------------
	//title := rows["title"]
	//if isWhiteSpaceString(title) {
	//	badRequest(w)
	//	return
	//}
	//price, err := toFloat(rows["price"])
	//if err != nil {
	//	internalServerError(w, err)
	//	return
	//}
	//thumbnail := rows["thumbnail"]
	//details := rows["details"]
	//specification := rows["specification"]
	//service := rows["service"]
	//properties := joinArray(rows["properties"])
	//showcases := joinArray(rows["showcases"])
	//// -----------------------------------
	//var uid string
	//err = e.DB.QueryRow(context.Background(), InsertStoreSQL,
	//	title,
	//	price,
	//	thumbnail,
	//	details,
	//	specification,
	//	service,
	//	properties,
	//	showcases).Scan(&uid);
	//if err != nil {
	//	internalServerError(w, err)
	//	return
	//}
	//// -----------------------------------
	//w.Write([]byte(uid))
}
