package handlers

import (
	"bytes"
	"fmt"
	"github.com/jackc/pgx/pgtype"
	"net/http"
	"store/common"
)

const (
	LikeSQL = "select * from store_list_like_results($1,$2,$3,$4)"
)

func ResultsHandler(e *common.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		keyword := r.URL.Query().Get("keyword")
		like := "%" + keyword + "%"
		items, err := e.DB.Fetch(LikeSQL, like, 20, 0, 1)
		if err != nil {
			internalServerError(w, err)
			return
		}
		var writer bytes.Buffer

		for _, i := range items {

			item := i.([]interface{})

			uid := item[0].(string)
			title := item [1].(string)
			price, err := item[2].(*pgtype.Numeric).Value()
			if err != nil {
				internalServerError(w, err)
				return
			}
			thumbnail := item [3].(string)
			quantities, ok := item[4].(int32)
			if !ok {
				quantities = 0
			}

			writer.WriteString(`<div class="skelecton-item"><a class="item-link" href="/store/details/`)
			writer.WriteString(uid)
			writer.WriteString(`"><div class="item-image"><img class="image_src" data-src="/store/static/pictures/`)
			writer.WriteString(thumbnail)
			writer.WriteString(`"/></div><div class="item-info"><div class="item-info_title"><span>`)
			writer.WriteString(title)
			writer.WriteString(`</span></div><div class="item-info_count"><div class="count_price"><i>￥</i>`)
			writer.WriteString(fmt.Sprintf("%.2f", float(price)))
			writer.WriteString(`</div><div class="count_vol">`)
			writer.WriteString(fmt.Sprintf("成交 %d 笔", quantities))
			writer.WriteString(`</div></div></div></a></div>`)
		}

		m := map[string]interface{}{
			"Title": fmt.Sprintf("搜索 \"%s\"", keyword),
			"Debug": e.Debug,
			"Items": writer.String(),
		}
		renderPage(w, "results.html", m, e.Debug)

	})
}
