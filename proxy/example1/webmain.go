package main

import (
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
)

type web1handler struct {
}

// 客户端应答
// 假设密码是root root
// 会把连着拼接成 root:root
// 然后base64变成       xxxxx
// 请求头是发送添加头:    Authorization:Basic xxxxx

func (web1handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	auth := request.Header.Get("Authorization")
	if auth == "" {

		//请求头包含内容 WWW-Authenticate:Basic realm="您必须输入用户名和密码"
		writer.Header().Set("WWW-Authenticate", `Basic realm="您必须输入用户名和密码"`)
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	// Authorization: Basic c2hlbnlpOjEyMw==
	auth_list := strings.Split(auth, " ")

	if len(auth_list) == 2 && auth_list[0] == "Basic" {
		res, err := base64.StdEncoding.DecodeString(auth_list[1])
		if err == nil && string(res) == "root:root" {
			writer.Write([]byte("<h1>web1</h1>"))
			return
		}
	}
	writer.Write([]byte("用户名密码错误"))

}

type web2handler struct{}

func (web2handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("web2"))
}

func main() {

	c := make(chan os.Signal)
	go (func() {
		http.ListenAndServe(":9091", web1handler{})
	})()

	go (func() {
		http.ListenAndServe(":9092", web2handler{})
	})()
	signal.Notify(c, os.Interrupt)
	s := <-c
	log.Println(s)
}
