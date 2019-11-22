package handlers

import (
	"bytes"
	"context"
	"euphoria/blackfriday"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgtype"
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
	Properties    string
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

		err := e.DB.QueryRow(context.Background(), "select * from store_fetch_details($1)", uid).Scan(
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
			UId:           uid,
			Title:         title,
			Price:         fmt.Sprintf("%.2f", float(priceStr)),
			Details:       string(blackfriday.Run([]byte(details))),
			Specification: specification, Service: service, Showcases: showcases,
			Properties: buildProperties(properties),
			Taobao:     taobao.String,
			Quantities: quantities.Int,
			Debug:      e.Debug,
		}, e.Debug)
	})
}
func buildProperties(values []string) string {
	j := len(values)
	if j == 0 {
		return ""
	}
	var writer bytes.Buffer

	for i := 0; i < j; i += 2 {
		if i+1 >= j {
			return writer.String()
		}
		name := values[i]
		value := values[i+1]
		writer.WriteString(`<div class="detail-attribute-item"><div class="detail-attribute-item-container"><div class="detail-attribute-item-name">`)
		writer.WriteString(name)
		writer.WriteString(`</div><div class="detail-attribute-item-value">`)
		writer.WriteString(value)
		writer.WriteString(`</div></div></div>`)

	}
	return writer.String()
}
