package util

import (
	"math/rand"
	"time"
)

type HttpServer struct {  //目标server类
	Host string
}
func NewHttpServer(host string) *HttpServer  {
	return &HttpServer{Host:host}
}
type LoadBalance struct { //负载均衡类
	Servers []*HttpServer

}

func NewLoadBalance() *LoadBalance {
   return &LoadBalance{Servers:make([]*HttpServer,0)}
}

func(this *LoadBalance) AddServer(server *HttpServer)  {
	 this.Servers=append(this.Servers,server)
}
func(this *LoadBalance) SelectByRand() *HttpServer { //随机算法

	rand.Seed(time.Now().UnixNano())
	index:=rand.Intn(len(this.Servers))
	return this.Servers[index]
}
