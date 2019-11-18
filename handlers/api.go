package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"store/common"
)

const (
	InsertStoreSQL = "select * from store_insert($1,$2,$3,$4,$5,$6,$7,$8)"
	UpdateStoreSQL = "select * from store_update($1,$2,$3,$4,$5,$6,$7,$8,$9)"
	FetchStoreSQL  = "select title,price,thumbnail,details,specification,service,taobao,wholesale,properties,showcases from commodities where uid = $1"
	ListStoreSQL   = "select * from store_list($1,$2)"
)

func ApiCategoryHandler(e *common.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		items, err := readData(e, r)
		if err != nil {
			internalServerError(w, err)
			return
		}
		t, err := e.DB.Exec("select * from store_category_insert($1)", joinArray(*items))
		if err != nil {
			internalServerError(w, err)
			return
		}
		writeCommandTag(t, w)
	})
}
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
func ApiSlideHandler(e *common.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		items, err := readData(e, r)
		if err != nil {
			internalServerError(w, err)
			return
		}
		t, err := e.DB.Exec("select * from store_slide_insert($1)", joinArray(*items))
		if err != nil {
			internalServerError(w, err)
			return
		}
		writeCommandTag(t, w)
	})
}
func ApiStoreHandler(e *common.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method := r.URL.Query().Get("method")
		if r.Method == "POST" {
			switch method {
			case "insert":
				insertStore(e, w, r)
				return
			case "update":
				updateStore(e, w, r)
				return
			}
		}
		notFound(w)
	})
}
func ApiSellHandler(e *common.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method := r.URL.Query().Get("method")
		if method == "POST" {
			switch method {
			case "insert":
				insertSell(e, w, r)
				return
			}
		}
		notFound(w)
	})
}

func insertSell(e *common.Env, w http.ResponseWriter, r *http.Request) {

	rows, err := readMapString(e, r)
	if err != nil {
		internalServerError(w, err)
		return
	}
	// -----------------------------------

	uid := rows["uid"]
	if isWhiteSpaceString(uid) {
		badRequest(w)
		return
	}

	taobao := rows["taobao"]
	wholesaler := rows["wholesaler"]
	quantities := rows["quantities"]

	fmt.Println(taobao, wholesaler, quantities)

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
func insertSearch(e *common.Env, w http.ResponseWriter, r *http.Request) {
	items, err := readData(e, r)
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
	err = e.DB.QueryRow(InsertStoreSQL,
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
func updateStore(e *common.Env, w http.ResponseWriter, r *http.Request) {
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
	uid := rows["uid"]
	if isWhiteSpaceString(uid) {
		badRequest(w)
		return
	}
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
	t, err := e.DB.Exec(UpdateStoreSQL,
		uid,
		title,
		price,
		thumbnail,
		details,
		specification,
		service,
		properties,
		showcases)
	if err != nil {
		internalServerError(w, err)
		return
	}
	// -----------------------------------
	writeCommandTag(t, w)
}
