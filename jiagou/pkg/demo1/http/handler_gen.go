package http

import (
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/yb19890724/go-study/jiagou/pkg/demo1/endpoint"
	"net/http"
)

// 初始化监听的路由和处理的端点绑定
func NewHTTPHandler(endpoints endpoint.Endpoints, options map[string][]kithttp.ServerOption) http.Handler {
	m := mux.NewRouter()
	makeCreateHandler(m, endpoints, options["Create"])

	return m
}
