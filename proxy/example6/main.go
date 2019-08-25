package main

import (
	"fmt"
	"hash/crc32"
	"math/rand"
	"time"
)

type LoadBalance struct {
	server      []*httpServer
	ServerIndices []int
	CurIndex      int
}

type httpServer struct {
	Host   string
	Weight int
}

func (l *LoadBalance) AddService(hs *httpServer) {
	l.server = append(l.server, hs)
}

// 随机算法
func (l *LoadBalance) SelectByRand() *httpServer {
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(l.server))
	return l.server[index]
}

// 轮询算法
func (l *LoadBalance) RoundRobin() *httpServer {
	server := l.server[l.CurIndex]
	l.CurIndex = (l.CurIndex + 1) % len(l.server)
	return server
}

// ip_hash算法
func (l *LoadBalance) SelectByIPHash(ip string) *httpServer {
	index := int(crc32.ChecksumIEEE([]byte(ip))) % len(l.server)
	return l.server[index]
}

// 加权随机算法
func (l *LoadBalance) SelectByWeightRand() *httpServer {
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(l.ServerIndices))
	return l.server[l.ServerIndices[index]]
}

// 设置权重
func (l *LoadBalance) SetServiceWeight() {
	for index, v := range l.server {
		for i := 0; i <= v.Weight; i++ {
			l.ServerIndices = append(l.ServerIndices, index)
		}
	}
}

// 三个区间
// [0,3) , [3,5) [5,6)
// 全部区间[0,6)
func(l *LoadBalance) SelectByWeightRand2() *httpServer { // 加权随机算法(改良算法)
	
	var weight int
	rand.Seed(time.Now().UnixNano())
	sum:=0
	for i:=0;i<len(l.server);i++{
		sum+=l.server[i].Weight
	}
	rad:=rand.Intn(sum)
	for index, v :=range l.server {
		weight+=v.Weight
		if rad < weight {
			return l.server[index]
		}
	}
	return  l.server[0]
}

func main() {
	
	var (
		lb LoadBalance
	)
	
	lb.AddService(&httpServer{"127.0.0.1:8080", 3})
	lb.AddService(&httpServer{"127.0.0.2:8080", 2})
	lb.AddService(&httpServer{"127.0.0.3:8080", 1})
	
	lb.SetServiceWeight()
	
	for i := 0; i < 10; i++ {
		fmt.Println(lb.SelectByWeightRand2())
	}
}
