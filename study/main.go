package main

import (
	"fmt"
	"math/rand"
	"time"
)

type LoadBalance struct {
	Services      []*httpServer
	ServerIndices []int
	CurIndex      int
}

type httpServer struct {
	Host   string
	Weight int
}

func (l *LoadBalance) AddService(hs *httpServer) {
	l.Services = append(l.Services, hs)
}


// 设置权重
func (l *LoadBalance) SetServiceWeight() {
	for index, v := range l.Services {
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
	for i:=0;i<len(l.Services);i++{
		sum+=l.Services[i].Weight
	}
	rad:=rand.Intn(sum)
	
	
	for index, v :=range l.Services {
		weight+=v.Weight
		if rad < weight {
			return l.Services[index]
		}
	}
	return  l.Services[0]
}

func main() {
	
	var (
		lb LoadBalance
	)
	
	lb.AddService(&httpServer{"127.0.0.1:8080", 3})
	lb.AddService(&httpServer{"127.0.0.2:8080", 2})
	lb.AddService(&httpServer{"127.0.0.3:8080", 1})
	
	lb.SetServiceWeight()
	
	for i:=1;i<=10 ;i++  {
		fmt.Println(lb.SelectByWeightRand2())
	}
}
