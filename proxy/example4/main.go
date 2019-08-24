package main

import (
	"fmt"
	"math/rand"
	"time"
)

type HttpServer struct {  //目标server类
	Host string
	Weight int
}

func NewHttpServer(host string,weight int) *HttpServer  {
	return &HttpServer{Host:host,Weight:weight}
}

type LoadBalance struct { //负载均衡类
	Servers []*HttpServer
}

func(this *LoadBalance) AddServer(server *HttpServer)  {
	this.Servers=append(this.Servers,server)
}

// 加权随机算法
func main()  {
	
	var Lb LoadBalance
	
	rand.Seed(time.Now().UnixNano())
	
	// [0,5) [5,7) [7,8)
	// 在区间做筛选 命中那个区间 就调用那个服务
	
	
	Lb.AddServer(NewHttpServer("127.0.0.1",5))
	Lb.AddServer(NewHttpServer("127.0.0.2",2))
	Lb.AddServer(NewHttpServer("127.0.0.3",1))
	
	sumLimit :=make([]int,len(Lb.Servers))
	sum:=0
	for i:=0;i <len(Lb.Servers);i++{
		sum+=Lb.Servers[i].Weight
		sumLimit[i]=sum
	}
	
	
	rad :=rand.Intn(sum)
	fmt.Println(rad)
	
	
	fmt.Println(sumLimit)
	
	for index,value:=range sumLimit{
		
		if rad < value{
			fmt.Println(Lb.Servers[index])
		}
		
		fmt.Println(Lb.Servers[0])
	}
}
