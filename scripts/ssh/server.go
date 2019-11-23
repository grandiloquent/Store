package main

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

var (
	HostName string
	Password string
	UserName string
)

// ==============================================
// 格式化文件尺寸
var decimapAbbrs = []string{"B", "kB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}

func getSizeAndUnit(size float64, base float64, _map []string) (float64, string) {
	i := 0
	unitsLimit := len(_map) - 1
	for size >= base && i < unitsLimit {
		size = size / base
		i++
	}
	return size, _map[i]
}
func HumanSizeWithPrecision(size float64, precision int) string {
	size, unit := getSizeAndUnit(size, 1000.0, decimapAbbrs)
	return fmt.Sprintf("%.*g%s", precision, size, unit)
}
func HumanSize(size float64) string {
	return HumanSizeWithPrecision(size, 4)
}

// ==============================================
// 排序

type FileInfos []os.FileInfo

func (f FileInfos) Len() int {
	return len(f)
}
func (f FileInfos) Less(i, j int) bool {
	if f[i].IsDir() && f[j].IsDir() {
		return f[i].Name() < f[j].Name()
	} else if !f[i].IsDir() && !f[j].IsDir() {
		return f[i].Size() > f[j].Size()
	} else if f[i].IsDir() && !f[j].IsDir() {
		return true
	} else if !f[i].IsDir() && f[j].IsDir() {
		return false
	}
	return true
}
func (f FileInfos) Swap(i, j int) {
	tmp := f[i];
	f[i] = f[j];
	f[j] = tmp;
}

// ==============================================

func connect(user, password, host string, port int) (*ssh.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		err          error
	)
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))
	clientConfig = &ssh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: 30 * time.Second, HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	addr = fmt.Sprintf("%s:%d", host, port)
	if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}
	return sshClient, nil
}
func createDirectoryIfNotExists(client *sftp.Client, dir string) error {
	_, err := client.ReadDir(dir)
	if os.IsNotExist(err) {
		return client.MkdirAll(dir)
	}
	return err
}
func listDirectory(client *sftp.Client, p string) {
	fs, err := client.ReadDir(p)
	if err != nil {
		log.Fatal(err)
	}
	sort.Sort(FileInfos(fs))
	for _, f := range fs {
		if !f.IsDir() {
			fmt.Printf("% 8s %s\n", HumanSize(float64(f.Size())), f.Name());
		} else {
			fmt.Printf("%s\n", f.Name());
		}
	}
}
func loadingSettings() map[string]interface{} {
	buf, err := ioutil.ReadFile("../../settings/settings.json")
	if err != nil {
		log.Fatal(err)
	}
	var m map[string]interface{}
	err = json.Unmarshal(buf, &m)
	if err != nil {
		log.Fatal(err)
	}
	return m
}
func runCommand(sshClient *ssh.Client, cmd string) {
	fmt.Printf("Excuting the command -> %s.\n", cmd)
	session3, err := sshClient.NewSession()
	defer session3.Close()
	out, err := session3.CombinedOutput(cmd)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(out))
	}
}
func uploadFile(sftpClient *sftp.Client, localFile, remoteDir string) {
	// 用来测试的本地文件路径 和 远程机器上的文件夹
	srcFile, err := os.Open(localFile)
	if err != nil {
		log.Fatalf("%s. localFile -> %s \n", err.Error(), localFile)
	}
	defer srcFile.Close()
	var remoteFileName = filepath.Base(localFile)
	dstFile, err := sftpClient.Create(path.Join(remoteDir, remoteFileName))
	if err != nil {
		log.Fatalf("%s. localFile -> %s;  remoteDir -> %s\n", err.Error(), localFile, remoteDir)
	}
	defer dstFile.Close()
	buf := make([]byte, 8192)
	t := 0
	start := time.Now()
	total := 0
	startSize := 0
	for {
		n, _ := srcFile.Read(buf)
		if n == 0 {
			break
		}
		t = t + n
		dstFile.Write(buf[0:n])
		if time.Since(start).Seconds() > 1 {
			total++
			fmt.Printf("\n% s/秒 %d秒", HumanSize(float64(t-startSize)), total)
			start = time.Now()
			startSize = t;
		}
	}
	fmt.Printf("\n %s \n", remoteFileName)
}
func uploadDirectory(sftpClient *sftp.Client, localDirectory, remoteDir string, filter ( func(string) bool), allDirectory bool) {
	if allDirectory {
		filepath.Walk(localDirectory, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() && filter(path) {
				dir := remoteDir + strings.ReplaceAll(filepath.Dir(path)[len(localDirectory):], "\\", "/")
				createDirectoryIfNotExists(sftpClient, dir)
				uploadFile(sftpClient, path, dir)
			}
			return nil
		})
	} else {
		files, err := ioutil.ReadDir(localDirectory)
		if err != nil {
			log.Fatal(err)
		}
		for _, f := range files {
			if !f.IsDir() && filter(f.Name()) {
				uploadFile(sftpClient, localDirectory+"/"+f.Name(), remoteDir)
			}
		}
	}

}

// ==============================================

func cleanup() {
	// -----------------------------------
	c, f := connectSftp()
	defer c.Close()
	defer f.Close()
	// -----------------------------------
	runCommand(c, "systemctl stop GoCommodities")
	// -----------------------------------
	listDirectory(f, "/usr/bin")
	// -----------------------------------
	runCommand(c, "rm -f /usr/bin/commodities_linux_amd64")
}
func connectSftp() (*ssh.Client, *sftp.Client) {
	c, err := connect(UserName, Password, HostName, 22)
	if err != nil {
		log.Fatal(err)
	}
	f, err := sftp.NewClient(c)
	if err != nil {
		log.Fatal(err)
	}
	return c, f
}
func main() {
	setupEnv()
	// cleanup()
	//publishApplication()
	//publishService()
	//publishNginx()
	//publishFiles()

	publishTemplates()
	publishScript()
}
func publishApplication() {
	serverRoot := "/usr/bin"
	fileName := "../../store_linux_amd64"
	// /usr/bin/store_linux_amd64
	// -----------------------------------
	c, f := connectSftp()
	defer c.Close()
	defer f.Close()
	// -----------------------------------
	uploadFile(f, fileName, serverRoot)
}
func publishFiles() {
	serverDirectory := "/root/store"
	// -----------------------------------
	c, f := connectSftp()
	defer c.Close()
	defer f.Close()
	// -----------------------------------

	//settingsDirectory := path.Join(serverDirectory, "settings")
	//createDirectoryIfNotExists(f, settingsDirectory)
	//uploadFile(f, "../../settings/settings.json", settingsDirectory)

	// -----------------------------------

	//templateDirectory := path.Join(serverDirectory, "templates")
	//createDirectoryIfNotExists(f, templateDirectory)
	//uploadDirectory(f, "../../templates", templateDirectory, func(s string) bool {
	//	if strings.HasSuffix(s, ".html") {
	//		return true
	//	}
	//	return false
	//}, true)

	// -----------------------------------

	//runCommand(c, "rm -Rf "+serverDirectory+"/static")
	//
	staticDirectory := path.Join(serverDirectory, "static")
	createDirectoryIfNotExists(f, staticDirectory)
	uploadDirectory(f, "../../static", staticDirectory, func(s string) bool {
		if strings.HasSuffix(s, "js") || strings.HasSuffix(s, "css") {
			return true
		}
		return false
	}, false)

	// -----------------------------------

	//staticDirectory := path.Join(serverDirectory, "static")
	//createDirectoryIfNotExists(f, staticDirectory)
	//uploadDirectory(f, "../../static", staticDirectory, func(s string) bool {
	//	if strings.Contains(s,"\\pictures\\") {
	//		return true
	//	}
	//	return false
	//}, true)
}
func publishScript() {
	serverDirectory := "/root/store"
	// -----------------------------------
	c, f := connectSftp()
	defer c.Close()
	defer f.Close()
	// -----------------------------------

	// -----------------------------------

	runCommand(c, "rm -f "+serverDirectory+"/static/*.js")
	runCommand(c, "rm -f "+serverDirectory+"/static/*.css")

	staticDirectory := path.Join(serverDirectory, "static")
	createDirectoryIfNotExists(f, staticDirectory)
	uploadDirectory(f, "../../static", staticDirectory, func(s string) bool {
		if strings.HasSuffix(s, "js") || strings.HasSuffix(s, "css") {
			return true
		}
		return false
	}, false)


}
func publishTemplates() {
	serverDirectory := "/root/store"
	// -----------------------------------
	c, f := connectSftp()
	defer c.Close()
	defer f.Close()

	// -----------------------------------

	templateDirectory := path.Join(serverDirectory, "templates")
	createDirectoryIfNotExists(f, templateDirectory)
	uploadDirectory(f, "../../templates", templateDirectory, func(s string) bool {
		if strings.HasSuffix(s, ".html") {
			return true
		}
		return false
	}, true)

}
func publishNginx() {
	// -----------------------------------
	c, f := connectSftp()
	defer c.Close()
	defer f.Close()
	// -----------------------------------
	uploadFile(f, "nginx.conf", "/etc/nginx")
	runCommand(c, "systemctl restart nginx")
}
func publishService() {
	serviceName := "GoStore"
	// -----------------------------------
	c, f := connectSftp()
	defer c.Close()
	defer f.Close()
	// -----------------------------------
	uploadFile(f, serviceName+".service", "/etc/systemd/system")
	runCommand(c, "systemctl start "+serviceName)
	runCommand(c, "systemctl status "+serviceName)
	runCommand(c, "chmod +x /usr/bin/store_linux_amd64")
}
func setupEnv() {
	m := loadingSettings()
	server := m["Server"].(map[string]interface{})
	HostName = server["HostName"].(string)
	UserName = server["UserName"].(string)
	Password = server["Password"].(string)
}
