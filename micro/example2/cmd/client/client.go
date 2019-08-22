package main

import (
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/etcdv3"
	hello "github.com/yb19890724/go-study/micro/example2/proto"
	"time"

	"golang.org/x/net/context"
)

// 自定义实现熔断

// ----------------------------------//

type MyClientWrapper struct {
	client.Client
}

func (c *MyClientWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	return hystrix.Do(req.Service()+"."+req.Endpoint(), func() error {
		return c.Client.Call(ctx, req, rsp, opts...)
	}, func(e error) error {
		fmt.Println("服务降级")
		return e
	})
}

// NewClientWrapper returns a hystrix client Wrapper.
func NewClientWrapper() client.Wrapper {
	return func(c client.Client) client.Client {
		return &MyClientWrapper{c}
	}
}

// ----------------------------------//

type Say struct{}

func (g *Say) Hello(ctx context.Context, req *hello.Request, rsp *hello.Response) error {
	rsp.Msg = "Hello " + req.Name
	return nil
}

func main() {

	// hystrix.DefaultTimeout =4000

	reg := etcdv3.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{
			"http://127.0.0.1:2379",
		}
	})

	// 在此处创建新服务
	service := micro.NewService(
		micro.Name("hello.client"),
		micro.Registry(reg),
		// micro.Transport(grpc.NewTransport()),

		//micro.WrapClient(hys.NewClientWrapper()), // 熔断
		micro.WrapClient(NewClientWrapper()), // 自定义熔断 测试降级
		micro.WrapCall(),
	)
	service.Init() // 初始化

	cl := hello.NewSayService("hello.server", service.Client())

	t := time.NewTicker(100 * time.Millisecond)

	for e := range t.C {
		rsp, err := cl.Hello(context.TODO(), &hello.Request{Name: "Join"}, func(options *client.CallOptions) {
			// 重试可能和熔断冲突 请注释测试
			options.RequestTimeout = 10 * time.Second // 10秒
			options.Retry = func(ctx context.Context, req client.Request, retryCount int, err error) (b bool, e error) {

				fmt.Println("Retry")
				return false, nil
			}
		})

		if err != nil {
			fmt.Println(err)
		} else {
			// fmt.Printf("%v==%v",e,rsp.Msg)
			fmt.Println(rsp.Msg, e)
		}
		fmt.Println()

	}

}
