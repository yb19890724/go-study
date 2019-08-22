package main

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"golang.org/x/time/rate"
	"os"
)


type UserRequest struct {
	Name string
}

type UserServerer interface {
	GetUserName(name string) string
}

type UserServer struct {

}

func (u UserServer) GetUserName(name string) string {
	return name
}



type UserResponse struct {
	Result string `json:"result"`
}

// 日志中间件
func UserServerLoggerMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			logger.Log("method","get")
			
			return next(ctx,request)
		}
	}
}



// 限流判断  中间件 代码无侵入
func RateLimit(limit *rate.Limiter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if !limit.Allow(){
				return nil,errors.New("too many request")
			}
			return next(ctx,request)
		}
	}
}

// 后续正确逻辑
func getUserNameEndPoint (u UserServerer) endpoint.Endpoint{
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r :=request.(UserRequest)
		name:=u.GetUserName(r.Name)
		result:=name
		return  UserResponse{Result:result},nil
	}
}



func main() {
	
	var logger log.Logger
	{
		logger =log.NewLogfmtLogger(os.Stdout)
		logger =log.WithPrefix(logger,"my test","1.0")
		logger =log.With(logger,"time",log.DefaultTimestampUTC)
		logger =log.With(logger,"caller",log.DefaultCaller)
	}
	
	var limit = rate.NewLimiter(1, 5)
	user :=UserServer{}
	RateLimit(limit)(UserServerLoggerMiddleware(logger)(getUserNameEndPoint(user)))
	
}
