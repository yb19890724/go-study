package stringsvc3

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"golang.org/x/time/rate"

	"github.com/sony/gobreaker"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
	httptransport "github.com/go-kit/kit/transport/http"
)

// 根据用户输入的多个地址，创建到多个服务器的代理
func ProxyingMiddleware(ctx context.Context, instances string, logger log.Logger) ServiceMiddleware {
	// 如果是空不需要代理
	if instances == "" {
		logger.Log("proxy_to", "none")
		return func(next StringService) StringService { return next }
	}

	// 基础参数
	var (
		qps         = 100                    // 请求频率超过多少会返回错误
		maxAttempts = 3                      // 请求在放弃前重试多少次，用于 load balancer
		maxTime     = 250 * time.Millisecond // 请求在放弃前的超时时间，用于 load balancer
	)

	// 除此以外, 为列表中的每个实例构造一个端点，然后添加
	// 它到一组固定的端点。在真实的服务中，而不是这样做
	// 手动，您可能会使用package sd对您的服务的支持
	// 发现系统。
	var (
		instanceList = split(instances)
		endpointer   sd.FixedEndpointer
	)
	logger.Log("proxy_to", fmt.Sprint(instanceList))
	for _, instance := range instanceList {
		var e endpoint.Endpoint
		e = makeUppercaseProxy(ctx, instance)// 创建client
		e = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(e) //添加 breader
		e = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), qps))(e) //添加limiter
		endpointer = append(endpointer, e)
	}


	balancer := lb.NewRoundRobin(endpointer)  //添加load balancer
	
	// Retry封装一个service load balancer，返回面向特定service method的load balancer。到这个endpoint的请求会自动通过
	// load balancer进行分配到各个代理服务器中。返回失败的请求会自动retry直到成功或者到达最大失败次数或者超时。
	retry := lb.Retry(maxAttempts, maxTime, balancer)//添加retry机制

	// 最后，返回由proxymw 实现的ServiceMiddleware。
	return func(next StringService) StringService {
		return proxymw{ctx, next, retry}
	}
}

// 定义使用了代理机制的新服务
// proxymw实现StringService，将Uppercase请求转发给
// 提供端点，并通过服务器提供所有其他（即计数）请求
// next StringService.
type proxymw struct {
	ctx       context.Context
	next      StringService     // 通过此服务提供大多数请求..
	uppercase endpoint.Endpoint // ...除了uppercase，它由此端点提供服务
}

// 直接用当前服务处理Count请求
func (mw proxymw) Count(s string) int {
	return mw.next.Count(s)
}

// 将uppercase请求发往各个代理服务器中(后面会讲到通过Load balancer实现）

func (mw proxymw) Uppercase(s string) (string, error) {
	response, err := mw.uppercase(mw.ctx, uppercaseRequest{S: s})
	if err != nil {
		return "", err
	}

	resp := response.(uppercaseResponse)
	if resp.Err != "" {
		return resp.V, errors.New(resp.Err)
	}
	return resp.V, nil
}

// 创建到特定地址代理服务器的client
func makeUppercaseProxy(ctx context.Context, instance string) endpoint.Endpoint {
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	u, err := url.Parse(instance)
	if err != nil {
		panic(err)
	}
	if u.Path == "" {
		u.Path = "/uppercase"
	}
	return httptransport.NewClient(
		"GET",
		u,
		encodeRequest,
		decodeUppercaseResponse,
	).Endpoint()
}

// 获取用户指定的代理服务器地址列表,本样例中，用户输入多个代理服务器用”,”分割

func split(s string) []string {
	a := strings.Split(s, ",")
	for i := range a {
		a[i] = strings.TrimSpace(a[i])
	}
	return a
}
