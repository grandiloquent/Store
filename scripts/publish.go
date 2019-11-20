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

func main() {

	publishCss()
	publishJavaScript()
	publishBin()
}

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
		// strings.HasPrefix(v.Name(), "app") && (strings.HasSuffix(v.Name(), ".css") || strings.HasSuffix(v.Name(), ".js"))
		if strings.HasSuffix(v.Name(), ".css") || strings.HasSuffix(v.Name(), ".js") {
			os.Remove(dir + "/" + v.Name())
		}
	}
}

func substringBeforeLast(s string, sep uint8) string {
	for i := len(s) - 1; i > -1; i-- {
		if s[i] == sep {
			return s[0:i]
		}
	}
	return s
}

func publishJavaScript() {
	dir := "../templates"
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".html") {
			compressJavaScript(dir + "/" + f.Name())
		}
	}
}
func publishCss() {
	dir := flag.String("dir", "../static/css", "样式文件所在的目录")
	css := flag.String("css", "../templates/header.html", "引用样式文件的HTM文件路径")
	flag.Parse()
	deleteFiles("../static")
	mergeStyle(*css, *dir)

}

func publishBin() {
	cmdir := flag.String("cmdir", "C:/Users/psycho/go/src/store", "")
	cmd := exec.Command(*cmdir + "/linux.bat") // or whatever the program is
	cmd.Dir = * cmdir                          // or whatever directory it's in
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s", out);
	}
}
func compressJavaScript(path string) {
	dir := "../static"
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	regex := regexp.MustCompile("/js/([a-zA-Z0-9_]+).js\"></script>")
	matches := regex.FindAllStringSubmatch(string(buf), -1)
	var writer bytes.Buffer
	if len(matches) == 0 {
		return
	}
	for _, m := range matches {
		content := readString(dir + "/js/" + m[1] + ".js")
		writer.WriteString(content)
	}
	writerBuf := writer.Bytes()
	md := calculateMd5(&writerBuf)
	target := dir + "/" + strings.TrimSuffix(filepath.Base(path), ".html") + "." + md + ".js"
	ioutil.WriteFile(target, writerBuf, 0644)
	runCommand("uglifyjs", target, "-o", target)
	//regex = regexp.MustCompile("/static/([a-zA-Z]+\\.[a-zA-Z0-9_]{5,}).js\"></script>")
	regex = regexp.MustCompile("/static/[a-zA-Z0-9_.]+.js\"></script>")

	buf = regex.ReplaceAll(buf, []byte("/static/"+strings.TrimSuffix(filepath.Base(path), ".html")+"."+md+".js"+"\"></script>"))
	ioutil.WriteFile(path, buf, 0644)
}

func readString(path string) string {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return string(buf)

}

func mergeStyle(css, dir string) {
	buf, err := ioutil.ReadFile(css)
	if err != nil {
		log.Fatal(err)
	}
	regex := regexp.MustCompile("<link rel=\"stylesheet\" href=\"/store/static/css/([a-zA-Z0-9_\\-.]+)\"/>")
	find := regex.FindAllSubmatch(buf, -1)
	// -----------------------------------
	var w bytes.Buffer
	for _, f := range find {
		c := dir + "/" + string(f[1])
		buf, err = ioutil.ReadFile(c)
		if err != nil {
			log.Fatal(err)
		}
		w.Write(buf)
	}
	buf = w.Bytes()
	o := filepath.Join(substringBeforeLast(dir, '/'), "app."+calculate(&buf)+".css")
	err = ioutil.WriteFile(o, w.Bytes(), 0644)
	if err != nil {
		log.Fatal(err)
	}
	compressCSS(o)
	regex= regexp.MustCompile("rel=\"stylesheet\" href=\"/store/static/[a-zA-Z0-9._]+.css\"/>")

	buf, err = ioutil.ReadFile(css)
	if err != nil {
		log.Fatal(err)
	}
	buf = regex.ReplaceAll(buf, []byte("rel=\"stylesheet\" href=\"/store/static/"+strings.TrimSuffix(filepath.Base(o), ".css")+".css"+"\"/>"))
	ioutil.WriteFile(css, buf, 0644)
}
func calculateMd5(buf *[]byte) string {
	s := md5.Sum(*buf)
	m := make([]byte, 16)
	for i, b := range s {
		m[i] = b
	}
	return hex.EncodeToString(m)
}
func runCommand(name string, arg ...string) {
	cmd, err := exec.Command(name, arg...).CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(name, arg, string(cmd))
}
