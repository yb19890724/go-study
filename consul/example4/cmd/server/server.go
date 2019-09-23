package main

import (
	"context"
	"fmt"
	"github.com/yb19890724/go-study/consul/example4/pkg/consul"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/health/grpc_health_v1"
	"log"
	"net"
	// "server/internal/consul"
	// pb "server/proto/helloworld"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.Name)
	return &pb.HelloReply{Message: "Hello " + in.Name+port}, nil
}


func RegisterToConsul() {
	consul.RegitserService("127.0.0.1:8500", &consul.ConsulService{
		Name: "helloworld",
		Tag:  []string{"helloworld"},
		IP:   "127.0.0.1",
		Port: 50051,
	})
}

// 实现健康检查
type healthCheck struct {

}


// Check 实现健康检查接口，这里直接返回健康状态，这里也可以有更复杂的健康检查策略，比如根据服务器负载来返回
func (h *healthCheck) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,// 返回健康标识
	}, nil
}

func (h *healthCheck)	Watch(req *grpc_health_v1.HealthCheckRequest,srv grpc_health_v1.Health_WatchServer) error  {
	return nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Println(port)
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	grpc_health_v1.RegisterHealthServer(s,&healthCheck{})
	RegisterToConsul()
	
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}