/*
* @Author: scottxiong
* @Date:   2020-11-28 06:19:58
* @Last Modified by:   scottxiong
* @Last Modified time: 2020-12-05 15:57:36
 */
package main

import (
	"encoding/json"
	"fmt"
	"github.com/scott-x/gutils/fs"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"
)

var (
	configfile      string
	download_folder string
	wg              = sync.WaitGroup{}
	HOME            string
)

type YTB struct {
	Download_folder string `json:"download_folder"`
	Task_position   string `json:"task_position"`
}

func init() {
	HOME, flag := os.LookupEnv("HOME")
	if !flag {
		log.Println("$HOME should be set")
		return
	}
	configfile = HOME + "/.ytb/youtube-dl.json"
}

func main() {
	if !isFileExist(configfile) {
		log.Println("configuration file " + configfile + " not found, could you please check?")
		return
	}
	download_folder, task_position := getConfig(configfile)

	if download_folder == "" {
		download_folder = HOME + "/Desktop"
	}

	if task_position == "" {
		log.Println("The field task_position should not be null")
		return
	}

	fmt.Println("---------------- parse configuration ----------------")
	fmt.Printf("download folder: %s\n", download_folder)
	fmt.Printf("task file: %s\n", task_position)
	//get lines
	content, err := ioutil.ReadFile(task_position)
	if err != nil {
		panic(err)
	}

	if strings.TrimSpace(string(content)) == "" {
		log.Println("No task is set in " + task_position)
		return
	}
	con := strings.Trim(string(content), "\n")
	lines := strings.Split(con, "\n")

	//set proxy
	setProxy()

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			wg.Add(1)
			go download(line, task_position, download_folder)
		}
	}

	wg.Wait()

}

func deleteFinishedUrl(task_position, line, url string) {
	bs, err := ioutil.ReadFile(task_position)
	if err != nil {
		panic(err)
	}
	content := strings.ReplaceAll(string(bs), line, "")

	err = ioutil.WriteFile(task_position, []byte(content), 0644)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(url + " is deleted from " + task_position)
}

func download(line, task_position, download_folder string) {
	defer wg.Done()

	var url string = ""
	var to string = ""

	line = strings.TrimSpace(line)

	if strings.Contains(line, " ") {
		arr := strings.Split(line, " ")
		url = arr[0]
		d1 := path.Join(download_folder, arr[len(arr)-1])
		if !fs.IsExist(d1) {
			log.Printf("creating folder: %s\n", d1)
			fs.CreateDirIfNotExist(d1)
		}
		to = d1 + "/%(title)s.%(ext)s"
	} else {
		url = line
		to = download_folder + "/%(title)s.%(ext)s"
	}

	fmt.Println("start downloading ===> " + url)
	// getFormat(url)
	cmd := exec.Command("youtube-dl", "-i", "-c", "-o", to, url)
	err := cmd.Run()
	if err != nil {
		log.Printf("download %s error:%s\n", url, err)
		return
	}

	fmt.Println(url + ` ===> 100% downloaded`)

	deleteFinishedUrl(task_position, line, url)

}

func isFileExist(name string) bool {
	if fi, err := os.Stat(name); err == nil {
		if fi.Mode().IsRegular() {
			return true
		}
	}
	return false
}

func setProxy() {
	fmt.Println("---------------- checking http/https proxy ----------------")
	_, flag := os.LookupEnv("http_proxy")
	if !flag {
		err := os.Setenv("http_proxy", "http://127.0.0.1:1024")
		if err != nil {
			log.Printf("set $http_proxy error:%s", err)
		}
		log.Printf("$http_proxy has been set to http://127.0.0.1:1024\n")
	}

	_, flag = os.LookupEnv("https_proxy")
	if !flag {
		err := os.Setenv("https_proxy", "http://127.0.0.1:1024")
		if err != nil {
			log.Printf("set $https_proxy error:%s", err)
		}
		log.Printf("$https_proxy has been set to http://127.0.0.1:1024\n")
	}

}

func getFormat(url string) string {
	exec.Command("youtube-dl", "--get-format", url).Run()
	return ""
}

func getConfig(configfile string) (string, string) {
	bs, err := ioutil.ReadFile(configfile)
	if err != nil {
		panic(err)
	}

	ytb := &YTB{}
	err = json.Unmarshal(bs, ytb)
	if err != nil {
		panic(err)
	}
	return ytb.Download_folder, ytb.Task_position
}
