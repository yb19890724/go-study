package main

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"os"
)


type Endpoint func(request UserRequest) (name string, err error)

type Middleware func(Endpoint) Endpoint


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

// 后续正确逻辑
func getUserNameEndPoint (u UserServerer) Endpoint{
	return func(request UserRequest) (name string, err error) {
		name =u.GetUserName(request.Name)
		return name ,nil
	}
}



// 日志中间件
func UserServerLoggerMiddleware(logger log.Logger) Middleware {
	return func(next Endpoint) Endpoint {
		return func(request UserRequest) (name string, err error) {
			logger.Log("method","get")
			fmt.Println("logger")
			return next(request)
		}
	}
}



func RateLimit(l string) Middleware {
	return func(next Endpoint) Endpoint {
		return func (request UserRequest) (name string, err error) {
			fmt.Println(l)
			return next(request)
		}
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
	
	user :=UserServer{}
	
	e:=RateLimit("rate limit")(UserServerLoggerMiddleware(logger)(getUserNameEndPoint(user)))
	
	
	re :=UserRequest{Name:"test"}
	
	e(re)
	
}
