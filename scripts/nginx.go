package main

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
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

var decimapAbbrs = []string{"B", "kB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}

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
// 格式化文件尺寸

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

func cmd(sshClient *ssh.Client, cmd string) {
	fmt.Printf("cmd -> %s.\n", cmd)
	session3, err := sshClient.NewSession()
	defer session3.Close()
	out, err := session3.CombinedOutput(cmd)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(out))
	}
}
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
func downloadFile(sftpClient *sftp.Client, remoteFileName, localFileName string, perm os.FileMode) error {

	source, err := sftpClient.OpenFile(remoteFileName, os.O_RDONLY)

	if err != nil {
		return err
	}
	stat, err := source.Stat()
	if err != nil {
		return err
	}
	defer source.Close()

	target, err := os.OpenFile(localFileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, perm)
	if err != nil {
		return err
	}
	n, err := io.Copy(target, io.LimitReader(source, stat.Size()))
	if err != nil {
		return err
	}
	if n != stat.Size() {
		return fmt.Errorf("sftp 传输文件不完整")
	}

	err = os.Chmod(localFileName, perm)
	if err != nil {
	}
	return nil
}

// ==============================================
func loadingSettings() map[string]interface{} {
	buf, err := ioutil.ReadFile("../settings/settings.json")
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

func main() {
	m := loadingSettings()
	server := m["Server"].(map[string]interface{})
	HostName = server["HostName"].(string)
	UserName = server["UserName"].(string)
	Password = server["Password"].(string)

	dir := "1"
	args := os.Args
	args = args[1:]
	if len(args) == 0 {
		return
	}
	c, err := connect(UserName, Password, HostName, 22)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range args {
		if v == "app" {
			uploadApplication(c)

		} else if v == "nginx" {
			uploadNginx(c)
		} else if v == "static" {
			uploadStatic(c)
		} else if v == "templates" {
		} else if v == "js" {
			uploadJs(c)
			uploadCss(c)
			uploadTemplates(c)

		} else if v == "images" {
			uploadImages(c, dir)
		} else if v == "css" {
		}
	}
	defer c.Close()
	//c, err := connect(UserName, Password, HostName, 22)
	//if err != nil {
	//	log.Fatal(err)
	//}
	////setupService(c)
	////uploadTemplates(c, "C:/Users/psycho/go/src/commodities/templates", "/root/commodities/templates");
	//uploadApplication(c)
}
func uploadTemplates(c *ssh.Client) {

	cf, err := sftp.NewClient(c)
	if err != nil {
		log.Fatal(err)
	}

	source := "C:/Users/psycho/go/src/commodities/templates"
	target := "/root/commodities/templates"

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			t := filepath.Join(target, strings.ReplaceAll(path[len(source):], "\\", "/"))
			dir := strings.ReplaceAll(filepath.Dir(t), "\\", "/")

			err := createDirectoryIfNotExists(cf, dir)
			if err != nil {
				return err
			}
			uploadFile(cf, path, dir)
		}
		return nil
	})

}
func uploadCss(c *ssh.Client) {
	cf, err := sftp.NewClient(c)
	if err != nil {
		log.Fatal(err)
	}

	source := "C:/Users/psycho/go/src/commodities/static/css"
	target := "/root/commodities/static/css/";
	fs, err := cf.ReadDir(target)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range fs {
		cf.Remove(target + "/" + f.Name())
		fmt.Println("Delete -> ", target+f.Name())
	}
	cmd(c, "systemctl stop GoCommodities")

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasPrefix(info.Name(), "commodities") {
			t := filepath.Join(target, strings.ReplaceAll(path[len(source):], "\\", "/"))
			dir := strings.ReplaceAll(filepath.Dir(t), "\\", "/")

			err := createDirectoryIfNotExists(cf, dir)
			if err != nil {
				return err
			}
			fmt.Println("Upload -> ", path, dir)
			uploadFile(cf, path, dir)
		}
		return nil
	})

	uploadFile(cf, "C:/Users/psycho/go/src/commodities/templates/header.html", "/root/commodities/templates/")
	cmd(c, "systemctl start GoCommodities")
}

func uploadStatic(c *ssh.Client) {
	cf, err := sftp.NewClient(c)
	if err != nil {
		log.Fatal(err)
	}
	//source := "C:/Users/psycho/go/src/commodities/static"
	source := "C:/Users/psycho/Desktop/43ea6d5ccf18ba174e0e2e10f5f2e3f4.ico"
	//target := "/root/commodities/static/";
	target := "/root/euphoria/"

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			t := filepath.Join(target, strings.ReplaceAll(path[len(source):], "\\", "/"))
			dir := strings.ReplaceAll(filepath.Dir(t), "\\", "/")

			err := createDirectoryIfNotExists(cf, dir)
			if err != nil {
				return err
			}
			uploadFile(cf, path, dir)
		}
		return nil
	})

}
func uploadJs(c *ssh.Client) {
	cf, err := sftp.NewClient(c)
	if err != nil {
		log.Fatal(err)
	}
	source := "C:/Users/psycho/go/src/commodities/static/js"
	target := "/root/commodities/static/js/";
	fs, err := cf.ReadDir(target)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range fs {
		if strings.HasPrefix(f.Name(), "app") && strings.HasSuffix(f.Name(), ".js") {
			cf.Remove(target + "/" + f.Name())
			fmt.Println("Delete -> ", target+f.Name())
		}

	}

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasPrefix(info.Name(), "app") {
			t := filepath.Join(target, strings.ReplaceAll(path[len(source):], "\\", "/"))
			dir := strings.ReplaceAll(filepath.Dir(t), "\\", "/")

			err := createDirectoryIfNotExists(cf, dir)
			if err != nil {
				return err
			}
			uploadFile(cf, path, dir)
		}
		return nil
	})

}
func uploadImages(c *ssh.Client, source string) {
	cf, err := sftp.NewClient(c)
	if err != nil {
		log.Fatal(err)
	}

	target := "/root/commodities/static/images/";
	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".jpg") {
			t := filepath.Join(target, strings.ReplaceAll(path[len(source):], "\\", "/"))
			dir := strings.ReplaceAll(filepath.Dir(t), "\\", "/")

			err := createDirectoryIfNotExists(cf, dir)
			if err != nil {
				return err
			}
			uploadFile(cf, path, dir)
			//if _, err = cf.Stat(dir + "/" + info.Name()); os.IsNotExist(err) {
			//
			//}
		}
		return nil
	})

}
func setupService(c *ssh.Client) {
	cf, err := sftp.NewClient(c)
	if err != nil {
		log.Fatal(err)
	}

	// ==============================================

	//listDirectory(cf, "/etc/systemd/system");
	//err = downloadFile(cf, "/etc/systemd/system/GoEuphoria.service", "./GoEuphoria.service", 0644)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//uploadFile(cf, "./GoCommodities.service", "/etc/systemd/system")

	// ==============================================

	cmd(c, "systemctl start GoCommodities")
	cmd(c, "systemctl status GoCommodities")
	// journalctl -u GoCommodities.service --since today
	// /usr/bin/commodities_linux_amd64
	// chmod +x /usr/bin/commodities_linux_amd64
	defer c.Close()
	defer cf.Close()

}
func uploadApplication(c *ssh.Client) error {

	cf, err := sftp.NewClient(c)
	if err != nil {
		log.Fatal(err)
	}

	file := "C:/Users/psycho/go/src/commodities/commodities_linux_amd64"
	fileinfo, err := os.Stat(file)
	if os.IsExist(err) {
		return err
	}
	fmt.Printf("%s %s\n", fileinfo.Name(), HumanSize(float64(fileinfo.Size())))

	cmd(c, "systemctl stop GoCommodities")
	uploadFile(cf, file, "/usr/bin")
	cmd(c, "systemctl start GoCommodities")

	defer cf.Close()
	return nil
}
func uploadNginx(c *ssh.Client) {
	cf, err := sftp.NewClient(c)
	if err != nil {
		log.Fatal(err)
	}

	// ==============================================
	//err = downloadFile(cf, "/etc/nginx/nginx.conf", "./nginx.conf", 0644)
	//if err != nil {
	//	log.Fatal(err)
	//}
	// ==============================================
	uploadFile(cf, "./nginx.conf", "/etc/nginx")
	cmd(c, "systemctl restart nginx")
	defer cf.Close()
}
