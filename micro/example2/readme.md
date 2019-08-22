#### micro结合etcd做服务注册 服务发现

> micro 默认使用 本地127.0.0.1：2379的etcd


#### 生成proto

> proto > protoc --micro_out=. --go_out=. greeter.proto 


#### 限流时可以设置一下是否等待参数 true false