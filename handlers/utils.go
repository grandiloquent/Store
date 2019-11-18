package handlers

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"github.com/jackc/pgx"
	"github.com/pkg/errors"
	"io"
	"log"
	"net/http"
	"store/common"
	"strconv"
	"strings"
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
func writeCommandTag(t pgx.CommandTag, w http.ResponseWriter) {
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
