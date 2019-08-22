#### demo

运行 

```gotemplate

目录 cmd/

// 启动服务
go run main.go

// 用户调用
go run server.go --run_client

// 模拟clien 端

go run client.go

```


- 2019/08/07 10:19:36 Transport [http] Listening on [::]:54903 `传输协议和端口`
- 2019/08/07 10:19:36 Broker [http] Connected to [::]:54904    `代理连接`
- 2019/08/07 10:19:36 Registry [mdns] Registering node: greeter-c9e6dd18-21aa-470c-9b38-7b1dd692d4c2 `服务注册和服务发现`