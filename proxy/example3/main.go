package main

import (
	"fmt"
	"hash/crc32"
)

func main()  {
	// ip hash 算法
	var service =[]string{
		"127.0.0.1",
		"127.0.0.2",
		"127.0.0.3",
	}
	
	index:=int(crc32.ChecksumIEEE([]byte("192.168.72.188"))) % len(service)
	index1:=int(crc32.ChecksumIEEE([]byte("192.168.72.144"))) % len(service)
	fmt.Println(service[index],service[index1])
}
