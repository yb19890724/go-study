package main

import (
	. "github.com/yb19890724/go-study/proxy/example1/util"
	"log"
	"net/http"
	"net/http/httputil"
	url2 "net/url"
)

type ProxyHandler struct {}
func(* ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err:=recover();err!=nil{
			w.WriteHeader(500)
			log.Println(err)
		}
	}()
	
   lb:=NewLoadBalance()
   lb.AddServer(NewHttpServer("http://localhost:9091"))
   lb.AddServer(NewHttpServer("http://localhost:9092"))
   url,_:=url2.Parse(lb.SelectByRand().Host)
   proxy:=httputil.NewSingleHostReverseProxy(url)
   proxy.ServeHTTP(w,r)


}


func main()  {
	//fmt.Println(base64.StdEncoding.EncodeToString([]byte("shenyi:123")))
	http.ListenAndServe(":8080",&ProxyHandler{})

}
