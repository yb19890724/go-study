package main

import (
	"fmt"
	"net/rpc"
)



func main() {
	cli, err := rpc.DialHTTP("tcp", "127.0.0.1:10086")
	
	if err != nil {
		fmt.Println("网络连接失败")
	}
	
	var pd string
	
	cli.Call("Test.GetInfo","10086",&pd)
	
	fmt.Println(pd)
}
