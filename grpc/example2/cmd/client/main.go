package main

import (
	"context"
	pb "github.com/yb19890724/go-study/grpc/example2/proto"
	"google.golang.org/grpc"
	"log"
)

func main() {
	// gRPC 服务器的地址
	addr := "0.0.0.0:2333"
	
	// 不使用认证建立连接
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("connect server error: %v", err)
	}
	defer conn.Close()
	
	// 创建 gRPC 客户端实例
	grpcClient := pb.NewUserServiceClient(conn)
	
	// 调用服务端的函数
	req := pb.UserRequest{ID: 2}
	resp, err := grpcClient.GetUserInfo(context.Background(), &req)
	if err != nil {
		log.Fatalf("recevie resp error: %v", err)
	}
	
	log.Printf("[RECEIVED RESPONSE]: %v\n", resp) // 输出响应
}