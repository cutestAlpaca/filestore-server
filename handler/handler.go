package handler

import (
	"filestore-server/meta"
	"filestore-server/util"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// UploadHandler 处理文件上传
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// 返回页面
		file, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, "server error")
			return
		}

		io.WriteString(w, string(file))
	} else if r.Method == "POST" {
		// 接收文件到本地
		file, head, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("Failed get data err: %s\n", err.Error())
			return
		}
		defer file.Close()

		fileMeta := meta.FileMeta{
			FileName: head.Filename,
			Location: "/tmp/" + head.Filename,
			UploadAt: time.Now().Format(time.DateTime),
		}

		newFile, err := os.Create(fileMeta.Location)
		if err != nil {
			fmt.Printf("Failed create data err: %s\n", err.Error())
			return
		}
		defer newFile.Close()

		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("Failed save data err: %s\n", err.Error())
			return
		}

		newFile.Seek(0, 0)
		fileMeta.FileSha1 = util.FileSha1(newFile)
		meta.UpdateFileMeta(fileMeta)

		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}
}

func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload finished!")
}