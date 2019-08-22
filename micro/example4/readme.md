#### micro 网关demo

api 暴露http接口 转化成果rpc 调用srv的rpc服务

#### 启动

```
go run srv/main.go

go run api/api.go 

micro api --handler=api

curl http://localhost:8080/greeter/say/hello?name=John
```