package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"store/common"
	"text/template"
)

type Home struct {
	Title          string
	Debug          bool
	SearchHolder   string
	SearchKeywords []string
	Slide          []interface{}
	Items          string
}

func HomeHandler(e *common.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		searchKeywords, err := fetchSearchKeywords(e)
		if err != nil {
			internalServerError(w, err)
			return
		}
		slide, err := e.DB.Fetch("select filename from store_slide")
		if err != nil {
			internalServerError(w, err)
			return
		}

		items, err := e.DB.Fetch(ListStoreSQL, 10, 0)

		var writer bytes.Buffer

		for i := 0; i < len(items); i += 2 {
			writer.WriteString(`<div class="like-row">`)
			item := items[i].([]interface{})
			uid := item[0].(string)
			title := item [1].(string)
			price := stringPrice(item[2])
			thumbnail := item [3].(string)
			quantities, ok := item[4].(int32)
			if !ok {
				quantities = 0
			}
			writer.WriteString(`<div class="like-cell" data-id="`)
			writer.WriteString(uid)
			writer.WriteString(`"><img src="/store/static/pictures/`)
			writer.WriteString(thumbnail)
			writer.WriteString(`"/><div class="like-cell-footer"><span>`)
			writer.WriteString(title)
			writer.WriteString(`</span>`)
			writer.WriteString(`<div class="like-cell-tags"></div>`)
			writer.WriteString(`<div class="like-cell-bottom"><span class="like-price">￥`)
			writer.WriteString(fmt.Sprintf("%.2f", float(price)))
			writer.WriteString(`</span> <span class="like-quantities">`)
			writer.WriteString(fmt.Sprintf("已出售 %d  件", quantities))
			writer.WriteString(`</span></div></div></div>`)

			if i+1 < len(items) {
				item = items[i+1].([]interface{})
				uid = item[0].(string)
				title = item [1].(string)
				price = stringPrice(item[2])
				if err != nil {
					internalServerError(w, err)
					return
				}
				thumbnail = item [3].(string)
				quantities, ok = item[4].(int32)
				if !ok {
					quantities = 0
				}
				writer.WriteString(`<div class="like-cell" data-id="`)
				writer.WriteString(uid)
				writer.WriteString(`"><img src="/store/static/pictures/`)
				writer.WriteString(thumbnail)
				writer.WriteString(`"/><div class="like-cell-footer"><span>`)
				writer.WriteString(title)
				writer.WriteString(`</span>`)
				writer.WriteString(`<div class="like-cell-tags"></div>`)
				writer.WriteString(`<div class="like-cell-bottom"><span class="like-price">￥`)
				writer.WriteString(fmt.Sprintf("%.2f", float(price)))
				writer.WriteString(`</span> <span class="like-quantities">`)
				writer.WriteString(fmt.Sprintf("已出售 %d  件", quantities))
				writer.WriteString(`</span></div></div></div>`)
			}

			writer.WriteString(`</div>`)
		}

		writeHome(w, &Home{
			Title:          "淘货",
			Debug:          e.Debug,
			SearchHolder:   "精选好货",
			SearchKeywords: searchKeywords,
			Slide:          slide,
			Items:          writer.String(),
		})
	})
}
func writeHome(w http.ResponseWriter, data *Home) {
	t, err := template.ParseFiles("templates/home.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		internalServerError(w, err)
		return
	}
	t.Execute(w, *data)
}
