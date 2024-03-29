package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"store/common"
)

const (
	InsertStoreSQL = "select * from store_insert($1,$2,$3,$4,$5,$6,$7,$8,$9)"
	UpdateStoreSQL = "select * from store_update($1,$2,$3,$4,$5,$6,$7,$8,$9)"
	FetchStoreSQL  = "select title,price,thumbnail,details,specification,service,properties,showcases from store where uid = $1"
	ListStoreSQL   = "select * from store_list($1,$2)"
	InsertSellSQL  = "select * from store_sell_insert($1,$2,$3,$4)"
)

func ApiCategoryHandler(e *common.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		items, err := readData(e, r)
		if err != nil {
			internalServerError(w, err)
			return
		}
		t, err := e.DB.Exec(context.Background(), "select * from store_category_insert($1)", joinArray(*items))
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
				insertSearchKeywords(e, w, r)
				return
			}
		} else {
			switch method {
			case "fetch":
				fetchSearch(e, w, r)
				return
			case "like":
				fetchLikeSearch(e, w, r)
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
		t, err := e.DB.Exec(context.Background(), "select * from store_slide_insert($1)", joinArray(*items))
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
				insertStores(e, w, r)
				return
			case "update":
				updateStore(e, w, r)
				return
			}
		} else if r.Method == "GET" {
			switch method {
			case "details":
				fetchStoreDetails(e, w, r)
				return
			default:
				fetchStoreList(e, w, r)
				return
			}
		}
		notFound(w)
	})
}
func ApiSellHandler(e *common.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method := r.URL.Query().Get("method")
		if r.Method == "POST" {
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
	quantities, ok := rows["quantities"].(float64)
	if !ok {
		badRequest(w)
		return
	}
	// -----------------------------------
	t, err := e.DB.Exec(context.Background(), InsertSellSQL, uid, taobao, wholesaler, int64(quantities))
	if err != nil {
		internalServerError(w, err)
		return
	}
	writeCommandTag(t, w)

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
	t, err := e.DB.Exec(context.Background(), UpdateStoreSQL,
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
func fetchStoreDetails(e *common.Env, w http.ResponseWriter, r *http.Request) {
	if !checkAuthorization(r, e.AccessToken) {
		forbidden(w)
		return
	}
	uid := r.URL.Query().Get("uid")

	buf, err := fetchStore(uid, e)
	if err != nil {
		internalServerError(w, err)
		return
	}
	writeJson(w, buf)
	// http://localhost:5050/store/api/store?method=details&uid=lpkxgi

}
func fetchStore(uid string, e *common.Env) ([]byte, error) {

	var (
		title         string
		price         float64
		thumbnail     string
		details       string
		specification string
		service       string
		properties    []string
		showcases     []string
	)

	err := e.DB.QueryRow(context.Background(), FetchStoreSQL, uid).Scan(
		&title,
		&price,
		&thumbnail,
		&details,
		&specification,
		&service,
		&properties,
		&showcases,

	)
	if err != nil {
		return nil, err
	}
	m := make(map[string]interface{})

	m["uid"] = uid
	m["title"] = title
	m["price"] = common.ToFixed(price, 2)
	m["thumbnail"] = thumbnail
	m["details"] = details
	m["specification"] = specification
	m["service"] = service
	m["properties"] = properties
	m["showcases"] = showcases
	buf, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
func fetchLikeSearch(e *common.Env, w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("keyword")
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")
	sorttype := r.URL.Query().Get("sorttype")

	items, err := e.DB.Fetch(LikeSQL, "%"+keyword+"%", limit, offset, sorttype)
	if err != nil {
		internalServerError(w, err)
		return
	}
	var results []interface{}

	for _, i := range items {

		item := i.([]interface{})

		m := make(map[string]interface{})
		uid := item[0].(string)
		title := item [1].(string)
		price := item[2].(float64)

		thumbnail := item [3].(string)
		quantities, ok := item[4].(int32)
		if !ok {
			quantities = 0
		}

		m["uid"] = uid
		m["title"] = title
		m["thumbnail"] = thumbnail
		m["price"] = fmt.Sprintf("%.2f", float(price))
		m["quantities"] = quantities

		results = append(results, m)
	}
	buf, err := json.Marshal(results)
	if err != nil {
		internalServerError(w, err)
		return
	}
	writeJson(w, buf)
}
