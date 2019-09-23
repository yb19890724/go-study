package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/rpc"
)


//
type Test struct {

}

func (t *Test) GetInfo(argType string ,replyType *string) error  {
	fmt.Println("打印请求内容",argType)
	
	*replyType =argType+"123123"
	
	return  nil
}

func getText(w http.ResponseWriter,r *http.Request)  {
	io.WriteString(w,"hello test")
}


func main() {
	http.HandleFunc("/test",getText)
	
	ln,err:=net.Listen("tcp",":10086")
	
	// 创建类
	t:=new(Test)
	
	rpc.Register(t) // 注册rpc服务
	
	rpc.HandleHTTP()// 连接网络
	
	if err !=nil {
		fmt.Println("网络错误")
	}
	
	http.Serve(ln,nil)
}
