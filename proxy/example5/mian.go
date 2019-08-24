package main

import (
	"fmt"
)


type HttpServer struct {//目标server类
	Host string
	Weight int
	Cweight int
}

func NewHttpServer(host string,weight int) *HttpServer  {
	return &HttpServer{Host:host,Weight:weight,Cweight:0}
}

type LoadBalance struct { //负载均衡类
	Servers []*HttpServer
	CurrentIndex int  // 当前访问
}

func(this *LoadBalance) AddServer(server *HttpServer)  {
	this.Servers=append(this.Servers,server)
}

func loadBalanceRobin()  {
	var Lb LoadBalance
	Lb.AddServer(NewHttpServer("127.0.0.1",5))
	Lb.AddServer(NewHttpServer("127.0.0.2",2))
	Lb.AddServer(NewHttpServer("127.0.0.3",1))
	
	server :=Lb.Servers[Lb.CurrentIndex]
	fmt.Println(server)
	Lb.CurrentIndex =(Lb.CurrentIndex+1) % len(Lb.Servers)
	fmt.Println(Lb.Servers[Lb.CurrentIndex])
	Lb.CurrentIndex =(Lb.CurrentIndex+1) % len(Lb.Servers)
	fmt.Println(Lb.Servers[Lb.CurrentIndex])
	Lb.CurrentIndex =(Lb.CurrentIndex+1) % len(Lb.Servers)
	fmt.Println(Lb.Servers[Lb.CurrentIndex])
	Lb.CurrentIndex =(Lb.CurrentIndex+1) % len(Lb.Servers)
	fmt.Println(Lb.Servers[Lb.CurrentIndex])
	Lb.CurrentIndex =(Lb.CurrentIndex+1) % len(Lb.Servers)
	fmt.Println(Lb.Servers[Lb.CurrentIndex])
	Lb.CurrentIndex =(Lb.CurrentIndex+1) % len(Lb.Servers)
	//1 2 3 -- 1 2 3一直轮询
}

var ServerIndices[]int

func setServiceWeight(lb LoadBalance) {
	for index, server := range lb.Servers {
		if server.Weight > 0 {
			for i := 0; i < server.Weight; i++ {
				ServerIndices = append(ServerIndices, index)
			}
		}
	}
	
}

// 轮询加权
func loadBalanceRobinWeight()  {
	var Lb LoadBalance
	Lb.AddServer(NewHttpServer("127.0.0.1",3))
	Lb.AddServer(NewHttpServer("127.0.0.2",1))
	Lb.AddServer(NewHttpServer("127.0.0.3",1))
	setServiceWeight(Lb)
	
	fmt.Println(ServerIndices)
	
	// [0 0 0 1 2] 记录输出位置的次数
	for i:=0;i<len(ServerIndices);i++ {
		server :=Lb.Servers[ServerIndices[Lb.CurrentIndex]]
		Lb.CurrentIndex =(Lb.CurrentIndex+1) % len(ServerIndices)
		fmt.Println(server)
	}
}

// 轮询
func main()  {
	//loadBalanceRobin()
	//loadBalanceRobinWeight()

	
}

