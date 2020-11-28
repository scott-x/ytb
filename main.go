/*
* @Author: scottxiong
* @Date:   2020-11-28 06:19:58
* @Last Modified by:   scottxiong
* @Last Modified time: 2020-11-28 06:56:28
 */
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
)

var (
	where string
	to    string
	wg    = sync.WaitGroup{}
	HOME  string
)

type YTB struct {
	To            string `json:"to"`
	Task_position string `json:"task_position"`
}

func init() {
	HOME, flag := os.LookupEnv("HOME")
	if !flag {
		log.Println("$HOME should be set")
		return
	}
	where = HOME + "/.ytb/youtube-dl.json"
}

func main() {
	if !isFileExist(where) {
		log.Println("configuration file " + where + " not found, could you please check")
		return
	}
	to, task_position := getConfig(where)

	if len(to) == 0 {
		to = HOME + "/Desktop/%(title)s.%(ext)s"
	}

	if task_position == "" {
		log.Println("The field task_position should not be null")
		return
	}

	//get urls
	content, err := ioutil.ReadFile(task_position)
	if err != nil {
		panic(err)
	}

	if strings.TrimSpace(string(content)) == "" {
		log.Println("No task is set in " + task_position)
		return
	}
	con := strings.Trim(string(content), "\n")
	urls := strings.Split(con, "\n")

	//set proxy
	setProxy()

	for _, url := range urls {
		wg.Add(1)
		// download(to, "https://www.youtube.com/watch?v=C6FhEonS-SU")
		go download(to, url, task_position)
	}

	wg.Wait()

}

func deleteFinishedUrl(task_position, url string) {
	bs, err := ioutil.ReadFile(task_position)
	if err != nil {
		panic(err)
	}
	content := strings.ReplaceAll(string(bs), url, "")
	err = ioutil.WriteFile(task_position, []byte(content), 0644)
	if err != nil {
		log.Fatal(err)
	}
	if len(url) > 0 {
		log.Println(url + " is deleted from " + task_position)
	}
}

func download(to, url, task_position string) {
	defer wg.Done()

	if len(url) > 0 {
		fmt.Println("start downloading ===> " + url)
	}

	cmd := exec.Command("youtube-dl", "-i", "-c", "-o", to, url)
	err := cmd.Run()
	if err != nil {
		log.Printf("download %s error:%s\n", url, err)
		return
	}
	if len(url) > 0 {
		fmt.Println(url + ` ===> 100% downloaded`)
	}

	deleteFinishedUrl(task_position, url)
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
	fmt.Println("set http/https proxy...")
	err := os.Setenv("http_proxy", "http://127.0.0.1:1024")
	if err != nil {
		log.Printf("set $http_proxy error:%s", err)
	}
	err = os.Setenv("https_proxy", "http://127.0.0.1:1024")
	if err != nil {
		log.Printf("set $https_proxy error:%s", err)
	}
}

func getConfig(where string) (string, string) {
	bs, err := ioutil.ReadFile(where)
	if err != nil {
		panic(err)
	}

	ytb := &YTB{}
	err = json.Unmarshal(bs, ytb)
	if err != nil {
		panic(err)
	}
	return ytb.To, ytb.Task_position
}
