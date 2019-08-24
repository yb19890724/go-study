package main

import (
	"fmt"
	"hash/crc32"
)

type LoadBalance struct {
	Services []*httpService
}

type httpService struct {
	Host string
}

func (l *LoadBalance) AddService(hs *httpService) {
	l.Services = append(l.Services, hs)
}

// ip_hash算法
func (l *LoadBalance) SelectByIPHash(ip string) *httpService {
	index := int(crc32.ChecksumIEEE([]byte(ip))) % len(l.Services)
	return l.Services[index]
}

func main() {
	
	var lb LoadBalance
	
	lb.AddService(&httpService{"127.0.0.1:8080"})
	lb.AddService(&httpService{"127.0.0.2:8080"})
	lb.AddService(&httpService{"127.0.0.3:8080"})
	
	// is_hash
	fmt.Println(lb.SelectByIPHash("192.168.1.1"))
	fmt.Println(lb.SelectByIPHash("192.168.1.6"))
	fmt.Println(lb.SelectByIPHash("192.168.1.9"))
}
