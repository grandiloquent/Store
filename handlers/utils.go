package handlers

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	"github.com/pkg/errors"
	"io"
	"log"
	"net/http"
	"store/common"
	"strconv"
	"strings"
	"text/template"
	"unicode"
)

func isWhiteSpaceString(i interface{}) bool {
	s, ok := i.(string)
	if !ok {
		return false
	}

	for _, v := range s {
		if !    unicode.IsSpace(v) {
			return false
		}
	}
	return true
}
func toFloat(i interface{}) (float64, error) {
	f, ok := i.(float64)
	if ok {
		return f, nil;
	}
	return 0, errors.New("invalid")
}
func float(i interface{}) (float64) {
	f, ok := i.(string)
	if ok {
		v, err := strconv.ParseFloat(f, 64)
		if err != nil {
			return 0
		}
		return v
	}
	return 0
}
func readMapString(e *common.Env, r *http.Request) (map[string]interface{}, error) {
	items, err := readData(e, r)
	if err != nil {
		return nil, err
	}
	rows, ok := (*items).(map[string]interface{})
	if !ok {

		return nil, errors.New("")
	}
	return rows, nil
}
func badRequest(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}
func checkAuthorization(r *http.Request, accessToken string) bool {
	a := r.Header.Get("Authorization")
	b := "Bearer "
	if len(b)+len(accessToken) == len(a) && strings.HasPrefix(a, b) && strings.HasSuffix(a, accessToken) {
		return true
	}
	return false
}
func forbidden(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
}
func internalServerError(w http.ResponseWriter, err error) {
	log.Printf("error, %s", err.Error())
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
func joinArray(i interface{}) string {
	if i == nil {
		return "{}"
	}
	a := i.([]interface{})
	if a == nil {
		return ""
	}
	var buf bytes.Buffer
	buf.WriteString("{")
	if len(a) > 0 {
		s, ok := a[0].(string)
		if ok {
			buf.WriteString(s)
		}
	}
	if len(a) > 1 {
		for _, v := range a[1:] {
			s, ok := v.(string)
			if ok {
				buf.WriteString(",")
				buf.WriteString(s)
			}
		}
	}
	buf.WriteString("}")
	return buf.String()
}
func notFound(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}
func readData(e *common.Env, r *http.Request) (*interface{}, error) {
	if !checkAuthorization(r, e.AccessToken) {
		return nil, errors.New("forbidden")
	}
	var items interface{}
	err := readJson(r, &items)
	if err != nil {
		return nil, err
	}
	return &items, nil
}
func readJson(r *http.Request, items *interface{}) error {
	defer r.Body.Close()
	var err error
	var reader io.ReadCloser
	switch r.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(r.Body)
		if err != nil {
			return err
		}
		defer reader.Close()
	default:
		reader = r.Body
	}
	decoder := json.NewDecoder(reader)
	return decoder.Decode(items)
}
func writeCommandTag(t pgconn.CommandTag, w http.ResponseWriter) {
	out := make(map[string]interface{})
	out["RowsAffected"] = t.RowsAffected()

	buf, err := json.Marshal(out)
	if err != nil {
		internalServerError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(buf)
}
func writeJson(w http.ResponseWriter, buf []byte) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(buf)
}

func renderPage(w http.ResponseWriter, filename string, data interface{}, debug bool) {
	t, err := template.ParseGlob("./templates/*.html")
	if err != nil {
		internalServerError(w, err)
		return
	}
	err = t.ExecuteTemplate(w, filename, data)
	if err != nil {
		log.Println(err)
	}
}
func safeString(i interface{}) interface{} {
	var s interface{}
	v, ok := i.(string)
	if ok {
		s = v
	} else {
		s = nil
	}
	return s
}
func stringPrice(i interface{}) string {
	v, ok := i.(*pgtype.Numeric)
	if !ok {
		return ""
	}
	s, err := v.Value()
	if err != nil {
		return ""
	}
	f, err := strconv.ParseFloat(s.(string), 64)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%.2f", f)
}

func safeQueryInt(r *http.Request, key string, defaultValue int) int {
	v := r.URL.Query().Get(key)
	if v == "" {
		return defaultValue
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return defaultValue
	}
	return i
}
