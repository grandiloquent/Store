package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func calculate(buf *[]byte) string {
	s := md5.Sum(*buf)
	m := make([]byte, 16)
	for i, b := range s {
		m[i] = b
	}
	return hex.EncodeToString(m)
}
func checkFileType(filename string, extensions ...string) bool {
	l := strings.ToLower(filename)
	length := len(l)
	for length > 0 {
		length = length - 1
		if l[length] == '.' {
			l = l[length+1:]
			break
		}
	}
	for _, e := range extensions {
		if l == e {
			return true;
		}
	}
	return false
}
func compressCSS(path string) {
	runCommand("csso", path, "--output", path)
}
func deleteFiles(dir string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range files {
		if strings.HasPrefix(v.Name(), "app") && (strings.HasSuffix(v.Name(), ".css") || strings.HasSuffix(v.Name(), ".js")) {
			os.Remove(dir + "/" + v.Name())
		}
	}
}
func main() {
	dir := flag.String("dir", "../static/css", "样式文件所在的目录")
	css := flag.String("css", "../static/search-results.html", "引用样式文件的HTM文件路径")
	cmdir := flag.String("cmdir", "C:/Users/psycho/go/src/store", "")
	flag.Parse()
	deleteFiles(*dir)
	mergeStyle(*css, *dir)
	cmd := exec.Command(*cmdir + "/linux.bat") // or whatever the program is
	cmd.Dir = * cmdir                          // or whatever directory it's in
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s", out);
	}
}
func mergeStyle(css, dir string) {
	buf, err := ioutil.ReadFile(css)
	if err != nil {
		log.Fatal(err)
	}
	regex := regexp.MustCompile("<link rel=\"stylesheet\" href=\"([^\"]+)\"/>")
	find := regex.FindAllSubmatch(buf, -1)
	// -----------------------------------
	var w bytes.Buffer
	for _, f := range find {
		c := dir + "/" + string(bytes.Split(f[1], []byte("/"))[1])
		buf, err = ioutil.ReadFile(c)
		if err != nil {
			log.Fatal(err)
		}
		w.Write(buf)
	}
	buf = w.Bytes()
	o := filepath.Join(substringBeforeLast(dir, '/'), "app_"+calculate(&buf)+".css")
	err = ioutil.WriteFile(o, w.Bytes(), 0644)
	if err != nil {
		log.Fatal(err)
	}
	compressCSS(o)
}

func substringBeforeLast(s string, sep uint8) string {
	for i := len(s) - 1; i > -1; i-- {
		if s[i] == sep {
			return s[0:i]
		}
	}
	return s
}
