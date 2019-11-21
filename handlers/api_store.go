package handlers

import (
	"context"
	"github.com/jackc/pgx"
	"github.com/pkg/errors"
	"net/http"
	"store/common"
)

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

	rItems, err := readData(e, r)
	if err != nil {
		internalServerError(w, err)
		return
	}
	items, ok := (*rItems).([]interface{})
	if !ok {
		internalServerError(w, errors.New(""))
		return
	}

	batch := &pgx.Batch{}

	for _, v := range items {
		m, ok := v.(map[string]interface{})
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
		if !ok {
			internalServerError(w, errors.New(""))
			return
		}
		batch.Queue(InsertStoreSQL,
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