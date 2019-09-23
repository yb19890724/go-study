package main

import (
	"fmt"
	"hash/crc32"
	"math/rand"
	"sort"
	"time"
)

type HttpServers  []*HttpServer


func (p HttpServers) Len() int           { return len(p) }
func (p HttpServers) Less(i, j int) bool { return p[i].CWeight > p[j].CWeight } // 从大到小排序
func (p HttpServers) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type LoadBalance struct {
	Servers      HttpServers
	ServerIndices []int
	CurIndex      int
	SumWeight     int // 加权总和
}

type HttpServer struct {
	Host   string
	Weight int
	CWeight int // 当前权重
}

func (l *LoadBalance) AddService(hs *HttpServer) {
	l.Servers = append(l.Servers, hs)
}

// 随机算法
func (l *LoadBalance) SelectByRand() *HttpServer {
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(l.Servers))
	return l.Servers[index]
}

// 轮询算法
func (l *LoadBalance) RoundRobin() *HttpServer {
	server := l.Servers[l.CurIndex]
	l.CurIndex = (l.CurIndex + 1) % len(l.Servers)
	return server
}

// ip_hash算法
func (l *LoadBalance) SelectByIPHash(ip string) *HttpServer {
	index := int(crc32.ChecksumIEEE([]byte(ip))) % len(l.Servers)
	return l.Servers[index]
}

// 加权随机算法
func (l *LoadBalance) SelectByWeightRand() *HttpServer {
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(l.ServerIndices))
	return l.Servers[l.ServerIndices[index]]
}

// 设置权重
func (l *LoadBalance) SetServiceWeight() {
	for index, v := range l.Servers {
		for i := 0; i < v.Weight; i++ {
			l.ServerIndices = append(l.ServerIndices, index)
		}
	}
}

// 设置权重总和
func (l *LoadBalance) SetServiceSumWeight() {
	for _, v := range l.Servers {
		l.SumWeight=l.SumWeight+v.Weight
	}
}

// 三个区间
// [0,3) , [3,5) [5,6)
// 全部区间[0,6)
func(l *LoadBalance) SelectByWeightRand2() *HttpServer { // 加权随机算法(改良算法)
	
	var weight int
	rand.Seed(time.Now().UnixNano())
	sum:=0
	for i:=0;i<len(l.Servers);i++{
		sum+=l.Servers[i].Weight
	}
	rad:=rand.Intn(sum)
	for index, v :=range l.Servers {
		weight+=v.Weight
		if rad < weight {
			return l.Servers[index]
		}
	}
	return  l.Servers[0]
}

// 加权轮询
func(l *LoadBalance) RoundRobinByWeight() *HttpServer  { 
	server:=l.Servers[l.ServerIndices[l.CurIndex]]
	l.CurIndex=(l.CurIndex+1) % len(l.ServerIndices)
	return server
}

// 加权区间算法
func(l *LoadBalance) RoundRobinByWeight2() *HttpServer  {
	server:=l.Servers[0]
	sum:=0
	//3:1:1
	for i:=0;i<len(l.Servers);i++{
		sum+=l.Servers[i].Weight   //   第一次是3   [0,3)  [3,4)   [4,5)
		if l.CurIndex<sum {
			server=l.Servers[i]
			if l.CurIndex==sum-1 && i!=len(l.Servers)-1{
				l.CurIndex++
			} else {
				l.CurIndex=(l.CurIndex+1)%sum// 这里是重要的一步
			}
			fmt.Println(l.CurIndex)
			break
		}
	}
	return server
}


// 平衡加权轮询
func(l *LoadBalance) RoundRobinByWeight3() *HttpServer {
	for _,s:=range l.Servers{
		s.CWeight=s.CWeight+s.Weight
	}
	sort.Sort(l.Servers)
	max:=l.Servers[0] // 返回最大 作为命中服务
	
	max.CWeight=max.CWeight-l.SumWeight
	return max
}

func main() {
	
	var (
		lb LoadBalance
	)
	
	lb.AddService(&HttpServer{"127.0.0.1:8080", 3,0})
	lb.AddService(&HttpServer{"127.0.0.2:8080", 2,0})
	lb.AddService(&HttpServer{"127.0.0.3:8080", 1,0})
	
	lb.SetServiceWeight()
	lb.SetServiceSumWeight()
	
	for i := 0; i < 10; i++ {
		fmt.Println(lb.RoundRobinByWeight3())
	}
}
