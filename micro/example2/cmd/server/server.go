package main

import (
	"fmt"
	rl "github.com/juju/ratelimit"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/transport/grpc"
	"github.com/micro/go-plugins/registry/etcdv3"
	"github.com/micro/go-plugins/wrapper/ratelimiter/ratelimit"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/micro/go-micro"
	hello "github.com/yb19890724/go-study/micro/example2/proto"

	"golang.org/x/net/context"
)

var topic = "demo.topic"

type Say struct {
	Tag string
}

func (g *Say) Hello(ctx context.Context, req *hello.Request, rsp *hello.Response) error {
	// 测试超时  熔断和服务降级
	time.Sleep(10 * time.Second)
	rsp.Msg = "Hello " + req.Name + "Sever" + g.Tag
	return nil
}

func main() {

	reg := etcdv3.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{
			"http://127.0.0.1:2379",
		}
	})

	limit := 2 // 每秒钟两个请求
	b := rl.NewBucketWithRate(float64(limit), int64(limit))

	// 创建一个新服务 可选择包含一些选项。
	service := micro.NewService(
		micro.Name("hello.server"),
		micro.Registry(reg),
		micro.Version("0.0.1"),
		micro.Metadata(map[string]string{
			"type": "say",
		}),
		micro.Transport(grpc.NewTransport()),
		micro.WrapHandler(ratelimit.NewHandlerWrapper(b, false)), // 是否等待请求，true等待
	)

	service.Init()

	say := &Say{
		Tag: strconv.Itoa(rand.Int()),
	}

	fmt.Println("server tag:", say.Tag)
	err := hello.RegisterSayHandler(service.Server(), say)

	if err := broker.Init(); err != nil {
		log.Fatal(err)
	}

	if err := broker.Connect(); err != nil {
		log.Fatal(err)
	}

	go publisher()
	go subscribe()

	if err != nil {
		fmt.Println(err)
	}

	// run server
	if err := service.Run(); err != nil {
		panic(err)
	}
}

func publisher() {

	t := time.NewTicker(time.Second)
	for e := range t.C {
		msg := &broker.Message{
			Header: map[string]string{
				"Tag": strconv.Itoa(rand.Int()),
			},
			Body: []byte(e.String()),
		}

		if err := broker.Publish(topic, msg); err != nil {
			log.Fatal("[publish err]:$+v ", err)
		}
	}

}

func subscribe() {

	if _, err := broker.Subscribe(topic, func(event broker.Event) error {
		fmt.Printf("subscribe received msg : % s Header is %+v", event.Message().Body, event.Message().Header)
		fmt.Println()
		return nil
	}); err != nil {
		fmt.Printf("[subscribe err] : %+v", err)
	}

}
