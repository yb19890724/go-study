package service

import (
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/yb19890724/go-study/jiagou/pkg/demo1/endpoint"
	"github.com/yb19890724/go-study/jiagou/pkg/demo1/http"
	"github.com/yb19890724/go-study/jiagou/pkg/demo1/service"
	nethttp "net/http"
)

// 创建服务
func createService() nethttp.Handler {
	// 创建业务对象
	svc := service.New(nil)
	// 创建端点对象
	eps := endpoint.New(svc, nil)
	// 设置http服务服务中间件
	options := defaultHttpOptions()
	// 端点绑定到http服务上
	return http.NewHTTPHandler(eps, options)
}

// HTTP服务中间件（服务的aop）
func defaultHttpOptions() map[string][]kithttp.ServerOption {
	options := map[string][]kithttp.ServerOption{
		"Create": {
			kithttp.ServerErrorEncoder(http.ErrorEncoder),
		},
	}
	return options
}
