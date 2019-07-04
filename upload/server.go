package main
import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)
// 获取文件大小的接口
type Size interface {
	Size() int64
}
// 获取文件信息的接口
type Stat interface {
	Stat() (os.FileInfo, error)
}
// hello world, the web server
func HelloServer(w http.ResponseWriter, r *http.Request) {
	file, head, err := r.FormFile("file")
	
	data :=strings.Split(head.Header.Get("Content-Disposition"),";")
	fileType :=head.Header.Get("Content-Type")
	
	fmt.Fprintf(w, "上传文件名称: %s\r",head.Filename)
	fmt.Fprintf(w, "上传文件大小: %d\r",head.Size)
	fmt.Fprintf(w,"请求方式: %s\r",data[0])
	fmt.Fprintf(w,"文件名称: %s\r",data[1])
	fmt.Fprintf(w,"文件类型: %s\r",fileType)
	
	
	defer file.Close()
	
	
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if statInterface, ok := file.(Stat); ok {
		fileInfo, _ := statInterface.Stat()
		fmt.Fprintf(w, "上传文件的大小为: %d", fileInfo.Size())
	}
	if sizeInterface, ok := file.(Size); ok {
		fmt.Fprintf(w, "上传文件的大小为: %d", sizeInterface.Size())
	}
	
}
func main() {
	http.HandleFunc("/upload", HelloServer)
	err := http.ListenAndServe(":5050", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}