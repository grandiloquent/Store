package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func FileSystemExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	isExist := os.IsExist(err)
	return isExist
}

func DumpRequest(r *http.Request, outPath string) {
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile(outPath, dump, 0644)
	if err != nil {
		fmt.Println(err)
	}
}
func ParseQueryInt(u *url.URL, key string, defaultValue int) int {
	s := u.Query().Get(key)
	if s == "" || !IsDigit(s) {
		return defaultValue
	}
	c, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}
	return c
}
func FlatStrings(i []interface{}) []string {
	var s []string
	for _, v := range i {
		t, ok := v.([]interface{})
		if !ok || len(t) == 0 {
			continue
		}
		f, ok := t[0].(string)
		if !ok {
			continue
		}
		s = append(s, f)
	}
	return s
}
func DumpTypeName(v interface{}) {
	if v == nil {
		fmt.Println("nil")
	}
	fmt.Println(reflect.TypeOf(v).String())
}
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}
func InterfaceArray(i *interface{}) []interface{} {
	if i == nil {
		return nil
	}
	m, ok := (*i).([]interface{})
	if ok {
		return m
	}
	return nil
}
func IsDigit(s string) bool {
	if s == "" {
		return false
	}
	for _, v := range s {
		if !unicode.IsDigit(v) {
			return false
		}
	}
	return true
}
func IsLowerLetter(s string) bool {
	if s == "" {
		return false
	}
	for _, v := range s {
		if v < 'a' || v > 'z' {
			return false
		}
	}
	return true
}
func Join(separator string, contents []interface{}) string {
	if contents == nil || len(contents) == 0 {
		return ""
	}
	s := String(contents[0])
	if len(contents) > 1 {
		for i, v := range contents {
			if i == 0 {
				continue
			}
			ss := String(v)
			s = s + separator + ss
		}
	}
	return s
}
func MapString(i *interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	m, ok := (*i).(map[string]interface{})
	if ok {
		return m
	}
	return nil
}
func String(i interface{}) string {
	if i == nil {
		return ""
	}
	m, ok := i.(string)
	if ok {
		return m
	}
	return ""
}
func SubstringAfter(s string, c byte) string {
	p := -1
	j := len(s)
	for i := 0; i < j; i++ {
		if s[i] == c {
			p = i
			break
		}
	}
	if p == -1 {
		return s
	}
	return s[p+1:]
}
func SubstringAfterLast(s string, c byte) string {
	p := -1
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == c {
			p = i
			break
		}
	}
	if p == -1 {
		return s
	}
	return s[p+1:]
}
func Trim(s string) (bool, string) {
	if s == "" {
		return false, s
	}
	s = strings.TrimSpace(s)
	return s == "", s
}
func SubstringBeforeLast(s string, c byte) string {
	p := -1
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == c {
			p = i
			break
		}
	}
	if p == -1 {
		return s
	}
	return s[0:p]
}
func GetValidFileName(s string, separator string) string {

	reg := regexp.MustCompile("[\"<>|\\0:*?/\\\\]+")
	return reg.ReplaceAllString(s, separator)
}
