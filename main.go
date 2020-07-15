package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func initConfig(path string) map[string]string {
	// 读取key=value类型的配置文件
	config := make(map[string]string)

	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		panic(err)
	}

	r := bufio.NewReader(f)
	for {
		b, _, err := r.ReadLine()

		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		s := strings.TrimSpace(string(b))
		index := strings.Index(s, "=")
		if index < 0 {
			continue
		}
		key := strings.TrimSpace(s[:index])
		if len(key) == 0 {
			continue
		}
		value := strings.TrimSpace(s[index+1:])
		if len(value) == 0 {
			continue
		}
		config[key] = value

	}
	return config
}

var filenum int
var timestamp int64

func countFile(dir string) int {

	finfos, _ := ioutil.ReadDir(dir)
	for _, fi := range finfos {
		path := path.Join(dir, fi.Name())
		if fi.IsDir() {
			countFile(path)
			continue
		}
		osType := runtime.GOOS
		fileInfo, _ := os.Stat(path)
		if osType == "windows" {
			wFileSys := fileInfo.Sys().(*syscall.Win32FileAttributeData)
			tNanSeconds := wFileSys.CreationTime.Nanoseconds() /// 返回的是纳秒
			tSec := tNanSeconds / 1e9                          ///秒

			if timestamp > tSec {
				timestamp = tSec
			}
		}
		filenum++
	}
	return filenum
}

func main() {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	initconfigfileT := path.Join(dir, "monitor.conf")
	configs := initConfig(initconfigfileT)
	// configs := initConfig("monitor.conf")

	var dirs []string
	for i := 1; i > 0; i++ {
		dirname := "dir" + strconv.Itoa(i)
		configdir := configs[dirname]
		if len(configdir) > 0 {
			dirs = append(dirs, configdir)
		} else {
			break
		}

	}

	initconfigfile := path.Join(dir, "file_monitor.prom")
	file, _ := os.OpenFile(initconfigfile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)

	for i := 0; i < len(dirs); i++ {
		var currtime int64 = time.Now().Unix()
		timestamp = currtime
		errdir := os.Chdir(dirs[i])
		if errdir != nil {
			continue
		}
		num := countFile(dirs[i])
		files := strings.Split(dirs[i], "\\")
		dirname := files[cap(files)-1]
		dirname2 := files[cap(files)-2]
		if dirname == "" {
			dirname = dirname2
		}
		str1 := "node_file_count_nums{dirs=\"" + dirname + "\"} " + strconv.Itoa(num) + "\r\n"
		file.Write([]byte(str1)) //写入字节切片数据

		var difftime int64
		difftime = currtime - timestamp
		// fmt.Println(currtime, timestamp)

		str2 := "node_file_time_gap{dirs=\"" + dirname + "\"} " + strconv.FormatInt(difftime, 10) + "\r\n"
		file.WriteString(str2)
		filenum = 0
		timestamp = 0
	}
	file.Close()

	var paths = "C:\\Program Files\\wmi_exporter\\textfile_inputs\\file_monitor.prom"
	if configs["paths"] != "" {
		paths = configs["paths"]
	}

	err2 := os.Rename(initconfigfile, paths)
	fmt.Println(err2)

}
