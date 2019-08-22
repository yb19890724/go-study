package main

import (
	"fmt"
	"github.com/micro/go-micro"
	proto "github.com/yb19890724/go-study/mirco/example/proto"

	"golang.org/x/net/context"
)

type Greeter struct{}

func (g *Greeter) Hello(ctx context.Context, req *proto.HelloRequest, rsp *proto.HelloResponse) error {
	rsp.Greeting = "Hello " + req.Name
	return nil
}

func main() {
	// 在此处创建新服务
	service := micro.NewService(micro.Name("greeter.client"))
	service.Init() // 初始化

	greeter := proto.NewGreeterService("greeter", service.Client())
	rsp, err := greeter.Hello(context.TODO(), &proto.HelloRequest{Name: "test"})

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(rsp)
}
