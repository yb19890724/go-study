package main

import (
	"context"
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/yb19890724/go-study/consul/example2/proto/checks"
	pb "github.com/yb19890724/go-study/grpc/example2/proto"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

// 模拟的数据库查询结果
var users = map[int32]pb.UserResponse{
	1: {Name: "Dennis MacAlistair Ritchie", Age: 70},
	2: {Name: "Ken Thompson", Age: 75},
	3: {Name: "Rob Pike", Age: 62},
}

type simpleServer struct{}

// simpleServer 实现了 user.pb.go 中的 UserServiceServer 接口
func (s *simpleServer) GetUserInfo(ctx context.Context, req *pb.UserRequest) (resp *pb.UserResponse, err error) {
	if user, ok := users[req.ID]; ok {
		resp = &user
	}
	log.Printf("[RECEVIED REQUEST]: %v\n", req)
	return
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



// NewConsulRegister create a new consul register
func NewConsulRegister() *ConsulRegister {
	return &ConsulRegister{
		Address: "127.0.0.1:8500",
		Service: "unknown",
		Tag:     []string{"userinfo"},
		Port:    9527,
		DeregisterCriticalServiceAfter: time.Duration(1) * time.Minute,
		Interval:                       time.Duration(10) * time.Second,
	}
}

// ConsulRegister consul service register
type ConsulRegister struct {
	Address                        string
	Service                        string
	Tag                            []string
	Port                           int
	DeregisterCriticalServiceAfter time.Duration
	Interval                       time.Duration
}

// Register register service
func (r *ConsulRegister) Register() error {
	config := consulapi.DefaultConfig()
	config.Address = r.Address
	client, err := consulapi.NewClient(config)
	if err != nil {
		return err
	}
	agent := client.Agent()
	
	IP := localIP()
	reg := &consulapi.AgentServiceRegistration{
		ID:      fmt.Sprintf("%v-%v-%v", r.Service, IP, r.Port), // 服务节点的名称
		Name:    fmt.Sprintf("grpc.health.v1.%v", r.Service),    // 服务名称
		Tags:    r.Tag,                                          // tag，可以为空
		Port:    r.Port,                                         // 服务端口
		Address: IP,                                             // 服务 IP
		Check: &consulapi.AgentServiceCheck{ // 健康检查
			Name:fmt.Sprintf("grpc.health.v1.%v", r.Service),
			Interval: r.Interval.String(),                            // 健康检查间隔
			GRPC:     fmt.Sprintf("%v:%v/%v", IP, r.Port, r.Service), // grpc 支持，执行健康检查的地址，service 会传到 Health.Check 函数中
			DeregisterCriticalServiceAfter: r.DeregisterCriticalServiceAfter.String(), // 注销时间，相当于过期时间
		},
	}
	
	if err := agent.ServiceRegister(reg); err != nil {
		return err
	}
	
	return nil
}

func localIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}


func main() {
	
	// 指定微服务的服务端监听地址
	addr :=  localIP()+":9527"
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen error: %v", err)
	} else {
		log.Println("server listen: ", addr)
	}
	
	// 创建 gRPC 服务器实例
	grpcServer := grpc.NewServer()
	
	// 向 gRPC 服务器注册服务
	pb.RegisterUserServiceServer(grpcServer, &simpleServer{})
	
	
	r :=NewConsulRegister()
	r.Register()
	
	grpc_health_v1.RegisterHealthServer(grpcServer,&healthCheck{})
	
	// 启动 gRPC 服务器
	// 阻塞等待客户端的调用
	grpcServer.Serve(listener)
}
