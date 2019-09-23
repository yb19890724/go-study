#### consul健康检查api重写


1. protobuf写一个Health服务: 
https://github.com/grpc/grpc/blob/master/doc/health-checking.md

2. 检查名称
package必须写成grpc.health.v1 ... consul源代码里写死了调用grpc的health check



 protoc --go_out=plugins=grpc:. *.proto

#### 启动服务

过程：
1.注册grpc服务
2.注册consul rpc服务检查
3.localhost:8500/ui 查看是否成功