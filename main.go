/*
* @Author: scottxiong
* @Date:   2020-11-28 06:19:58
* @Last Modified by:   scottxiong
* @Last Modified time: 2020-11-28 06:56:28
 */
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

var (
	where string
	to    string
)

func init() {
	HOME, flag := os.LookupEnv("HOME")
	if !flag {
		log.Println("$HOME should be set")
		return
	}
	where = HOME + "/.ytb/youtube-dl.conf"
}

func main() {
	if !isFileExist(where) {
		log.Println("configuration file " + where + " not found, could you please check")
		return
	}

	bs, err := ioutil.ReadFile(where)
	if err != nil {
		panic(err)
	}

	if len(bs) == 0 {
		to = "~/Desktop/%(title)s.%(ext)s"
	} else {
		to = string(bs)
	}

	//set proxy
	setProxy()
	download(to, "https://www.youtube.com/watch?v=C6FhEonS-SU")

}

func downloadToWhere() string {
	return ""
}

func download(to, url string) {
	log.Println("start downloading " + url)
	cmd := exec.Command("youtube-dl", "-o", to, url)
	err := cmd.Start()
	if err != nil {
		log.Println("download " + url + " error")
	}
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
	fmt.Println("proxy")
}
