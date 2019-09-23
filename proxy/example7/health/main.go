package main

import (
	"fmt"
	"net/http"
	"time"
)

type LoadBalance struct {
	Servers       HttpServers
	ServerIndices []int
	CurIndex      int
}

type HttpServers []*HttpServer

type HttpServer struct {
	Host         string
	Status       string // 默认up 宕机：down
	FailCount    int    // 检测次数
	SuccessCount int    // 连续成功
}

func (l *LoadBalance) AddService(hs *HttpServer) {
	l.Servers = append(l.Servers, hs)
}

type HttpChecker struct {
	Servers      HttpServers
	FailMax      int // 被踢下线
 	RecoverCount int // 连续成功次数会被标记成 up
	FailFactor   float64 // 降权因子 默认是5.0
}

func NewHttpChecker(Servers HttpServers) *HttpChecker {
	return &HttpChecker{Servers: Servers, FailMax: 3, RecoverCount: 3}
}

func (h *HttpChecker) Check(time time.Duration) {
	client := http.Client{}
	for _, server := range h.Servers {
		res, err := client.Head(server.Host)
		if res != nil {
			defer res.Body.Close()
		}
		
		if err != nil { // 可能宕机
			h.Fail(server)
			continue
		}
		
		// 地址错误
		if res.StatusCode >= 200 && res.StatusCode < 400 {
			h.Success(server)
		} else {
			h.Fail(server)
		}
	}
	
}

func checkServers(Servers HttpServers) {
	// 定时器 设置时间
	t := time.NewTicker(time.Second * 3)
	check := NewHttpChecker(Servers)
	for {
		select {
		case <-t.C: // 根据设定时间往C中存入值
			check.Check(time.Second * 2)
			for _, s := range Servers {
				fmt.Println(s.Host, s.Status, s.FailCount)
			}
		}
	}
}

// 失败处理
func (h *HttpChecker) Fail(server *HttpServer) {
	if server.FailCount >= h.FailMax {
		server.Status = "DOWN"
	} else {
		
		server.FailCount++
	}
}

// 成功处理
func (h *HttpChecker) Success(server *HttpServer) {
	if server.FailCount > 0 {// 如果有失败打点
		server.FailCount--
		server.SuccessCount++
		
		// 达到连续成功次数 就设置Up
		if server.SuccessCount == h.RecoverCount {
			server.FailCount = 0
			server.Status = "UP"
			server.SuccessCount = 0
		}
	} else {
		server.Status = "UP"
	}
}

// 轮询算法
func (l *LoadBalance) RoundRobin() *HttpServer {
	server := l.Servers[l.CurIndex]
	l.CurIndex = (l.CurIndex + 1) % len(l.Servers)
	if server.Status == "DOWN" { // 宕机情况无法返回
		return l.RoundRobin()// 递归查找其他的服务
	}
	return server
}

func main() {
	
	var (
		lb LoadBalance
	)
	
	lb.AddService(&HttpServer{Host: "http://localhost:9091"})
	lb.AddService(&HttpServer{Host: "http://localhost:9092"})
	
	checkServers(lb.Servers)
	
}
