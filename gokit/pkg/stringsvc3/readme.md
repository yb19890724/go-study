#### 分析
- proxy 做反向代理
    - qps 限制
    - 尝试 次数
    - 最大响应时间
    
#### run三个main.go

- 8081
- 8080
- 8082
    这里加入代理
```go
        proxy  = flag.String("proxy", "localhost:8080,localhost:8081", "Optional comma-separated list of URLs to proxy uppercase requests")
```

