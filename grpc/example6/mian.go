package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/yb19890724/go-study/grpc/example6/protobuf"
)

func main(){
	text:=&test.Test{
		Name:"test",
		Tizhong:[]int32{150,155,168,180},
		Shengao:188,
		Motto:"就是干",
	}
	
	fmt.Println(text)
	
	data,err:=proto.Marshal(text)
	if err!=nil {
		fmt.Println("编码失败")
	}
	
	fmt.Println(data)
	
	newText:=&test.Test{}
	
	err =proto.Unmarshal(data,newText)
	
	if err!=nil {
		fmt.Println("解码失败")
	}
	
	fmt.Println(newText)
}
