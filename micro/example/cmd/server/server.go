package server

import (
	"fmt"
	"github.com/micro/cli"
	"os"

	"github.com/micro/go-micro"
	proto "github.com/yb19890724/go-study/mirco/example/proto"

	"golang.org/x/net/context"
)

type Greeter struct{}

func (g *Greeter) Hello(ctx context.Context, req *proto.HelloRequest, rsp *proto.HelloResponse) error {
	rsp.Greeting = "Hello " + req.Name
	return nil
}

// Setup and the client
func runClient(service micro.Service) {
	// Create new greeter client
	greeter := proto.NewGreeterService("greeter", service.Client())

	// Call the greeter
	rsp, err := greeter.Hello(context.TODO(), &proto.HelloRequest{Name: "John"})
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print response
	fmt.Println(rsp.Greeting)
}

func main() {
	// 创建一个新服务 可选择包含一些选项。
	service := micro.NewService(
		micro.Name("greeter"),
		micro.Version("latest"),
		micro.Metadata(map[string]string{
			"type": "helloworld",
		}),

		// 设置一些标志。指定--run_client以运行客户端

		// 添加运行时标志
		// 我们也可以在下面这样做
		micro.Flags(cli.BoolFlag{
			Name:  "run_client",
			Usage: "Launch the client",
		}),
	)

	// Init将解析命令行标志。任何标志设置将
	// 覆盖以上设置。这里定义的选项将
	// 覆盖命令行上设置的任何内容。
	service.Init(
		// 添加运行时操作
		// 我们实际上可以做到这一点
		micro.Action(func(c *cli.Context) {
			if c.Bool("run_client") {
				runClient(service)
				os.Exit(0)
			}
		}),
	)

	// 默认情况下，除非标志抓住，否则我们将运行服务器

	// 设置服务器

	// 注册处理程序
	proto.RegisterGreeterHandler(service.Server(), new(Greeter))

	// Run the server
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
