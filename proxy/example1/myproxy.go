package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

type ProxyHandler struct{}

func (*ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	
	defer func() {
		if err := recover(); err != nil {
			w.WriteHeader(500)
			log.Println(err)
		}
	}()
	if r.URL.Path == "/a" {
		newreq, _ := http.NewRequest(r.Method, "http://localhost:9091", r.Body)
		newresponse, _ := http.DefaultClient.Do(newreq)// 执行请求 返回响应
		defer newresponse.Body.Close() // 关闭连接
		res_cont, _ := ioutil.ReadAll(newresponse.Body)
		w.Write(res_cont)
		return
	}
	w.Write([]byte("default index"))
}

func main() {
	// fmt.Println(base64.StdEncoding.EncodeToString([]byte("shenyi:123")))
	http.ListenAndServe(":8080", &ProxyHandler{})
	
}
