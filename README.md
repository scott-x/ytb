# ytb

This tool is based on [youtube-dl](https://github.com/ytdl-org/youtube-dl), to use it please make sure golang was installed & Youtube can be visited.


### how to use?

1, downloading ytb tool

```bash
git clone git@github.com:scott-x/ytb.git
cd ytb
go run main.go
```

2, configuration

```bash
mkdir ~/.ytb
touch ~/.ytb/youtube-dl.json
```

`youtube-dl.json` has 2 properties as shown below:

```json
{
	"download_folder":"/Users/apple/Desktop", 
	"task_position":"/Users/apple/Desktop/task.txt" 
}
```

- `download_folder`: `download_folder` indicates the directory where these video will go to 
- `task_position`: before downloading, list the download youtube url line by line.


For Example:

![](./task.png)

```bash
~/go/src/github.com/scott-x/ytb(main*) » go run main.go                                       apple@iMac-52
---------------- parse configuration ----------------
download folder: /Users/apple/Desktop
task file: /Users/apple/Desktop/task.txt
---------------- checking http/https proxy ----------------
2020/12/05 13:02:48 $http_proxy has been set to http://127.0.0.1:1024
2020/12/05 13:02:48 $https_proxy has been set to http://127.0.0.1:1024
2020/12/05 13:02:48 creating folder: /Users/apple/Desktop/a/b
start downloading ===> https://www.youtube.com/watch?v=sBS22HtYEg8
start downloading ===> https://www.youtube.com/watch?v=Ttus8XGK6Xw
https://www.youtube.com/watch?v=sBS22HtYEg8 ===> 100% downloaded
2020/12/05 13:03:51 https://www.youtube.com/watch?v=sBS22HtYEg8 is deleted from /Users/apple/Desktop/task.txt
https://www.youtube.com/watch?v=Ttus8XGK6Xw ===> 100% downloaded
2020/12/05 13:04:32 https://www.youtube.com/watch?v=Ttus8XGK6Xw is deleted from /Users/apple/Desktop/task.txt
```
