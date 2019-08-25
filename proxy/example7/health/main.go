package main

import (
	"fmt"
	"net/http"
	"time"
)

type LoadBalance struct {
	servers      HttpServers
	ServerIndices []int
	CurIndex      int
}

type HttpServers []*httpServer

type httpServer struct {
	Host string
	Status string // 默认up 宕机：down
}

func (l *LoadBalance) AddService(hs *httpServer) {
	l.servers = append(l.servers, hs)
}


type HttpChecker struct {
	Servers HttpServers
}

func NewHttpChecker(servers HttpServers) *HttpChecker {
	return &HttpChecker{Servers:servers}
}

func (h *HttpChecker) Check(time time.Duration)  {
	client :=http.Client{}
	for _,server :=range h.Servers{
		res,err:=client.Head(server.Host)
		if res !=nil {
			defer res.Body.Close()
		}
		
		if err!=nil { // 可能宕机
			server.Status ="down"
		}
		
		// 地址错误
		if res.StatusCode >=200 && res.StatusCode<400 {
			server.Status ="up"
		}else {
			server.Status ="down"
		}
	}
	
}

func checkServers(servers HttpServers)  {
	// 定时器 设置时间
	t:=time.NewTicker(time.Second*3)
	check :=NewHttpChecker(servers)
	for {
		select {
		case <-t.C: // 根据设定时间往C中存入值
			check.Check(time.Second*2)
			for _,s :=range servers{
				fmt.Println(s.Host,s.Status)
			}
		}
	}
}

func main()  {
	
	var (
		lb LoadBalance
	)
	
	lb.AddService(&httpServer{Host:"http://localhost:9091"})
	lb.AddService(&httpServer{Host:"http://localhost:9092"})
	
		checkServers(lb.servers)
		
}
