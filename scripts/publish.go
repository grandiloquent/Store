package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
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
	regex = regexp.MustCompile("/static/([a-zA-Z]+\\.[a-zA-Z0-9_]{5,}).js\"></script>")
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
