package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/pgtype"
	"net/http"
	"store/common"
)

type Details struct {
	UId           string
	Title         string
	Price         string
	Details       string
	Specification string
	Service       string
	Showcases     []string
	Properties    []string
	Taobao        string
	Quantities    int32
	Debug         bool
}

func DetailsHandler(e *common.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		uid := vars["uid"]

		var title string
		var price pgtype.Numeric
		var details string
		var specification string
		var service string
		var showcases []string
		var properties []string
		var taobao pgtype.Text
		var quantities pgtype.Int4

		err := e.DB.QueryRow("select * from store_fetch_details($1)", uid).Scan(
			&title,
			&price,
			&details,
			&specification,
			&service,
			&showcases,
			&properties,
			&taobao,
			&quantities,
		)
		if err != nil {
			internalServerError(w, err)
			return
		}
		priceStr, err := price.Value()
		if err != nil {
			internalServerError(w, err)
			return
		}
		renderPage(w, "details.html", Details{
			UId:     uid,
			Title:   title,
			Price:   fmt.Sprintf("%.2f", float(priceStr)),
			Details: details, Specification: specification, Service: service, Showcases: showcases,
			Properties: properties,
			Taobao:     taobao.String,
			Quantities: quantities.Int,
			Debug:      e.Debug,
		}, e.Debug)
	})
}
