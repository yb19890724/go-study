package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
)
type web1handler struct {
}
func(web1handler) GetIP(request *http.Request) string{
	ips:=request.Header.Get("x-forwarded-for")
	if ips!=""{
	ips_list:= strings.Split(ips, ",")
		if len(ips_list)>0 && ips_list[0]!=""{
			return ips_list[0]
		}
	}
	return request.RemoteAddr
}


func(this web1handler) ServeHTTP(writer http.ResponseWriter, request *http.Request)  {
	writer.Write([]byte("web1"))
}
type web2handler struct {}

func(web2handler) ServeHTTP(writer http.ResponseWriter, request *http.Request)  {
	 writer.Write([]byte("web2"))
}

func main()  {
	c:=make(chan os.Signal)
	go(func() {
		http.ListenAndServe(":9091",web1handler{})
	})()
	go(func() {


		http.ListenAndServe(":9092",web2handler{})
	})()
	signal.Notify(c,os.Interrupt)
	s:=<-c
	log.Println(s)
}