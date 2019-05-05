package main

import (
	`github.com/yb19890724/go-study/gokit/pkg/stringsvc2`
	"net/http"
	"os"
	
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	
	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	httptransport "github.com/go-kit/kit/transport/http"
)

func main() {
	logger := log.NewLogfmtLogger(os.Stderr)
	
	fieldKeys := []string{"method", "error"}
	
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "my_group",
		Subsystem: "string_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "string_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)
	
	countResult := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "string_service",
		Name:      "count_result",
		Help:      "The result of each count method.",
	}, []string{}) // no fields here
	
	var svc stringsvc2.StringService
	svc = stringsvc2.Service{}
	// 日志记录
	svc = stringsvc2.LoggingMiddleware{logger, svc}
	// 数据收集
	svc = stringsvc2.InstrumentingMiddleware{requestCount, requestLatency, countResult, svc}
	
	uppercaseHandler := httptransport.NewServer(
		stringsvc2.MakeUppercaseEndpoint(svc),
		stringsvc2.DecodeUppercaseRequest,
		stringsvc2.EncodeResponse,
	)
	
	countHandler := httptransport.NewServer(
		stringsvc2.MakeCountEndpoint(svc),
		stringsvc2.DecodeCountRequest,
		stringsvc2.EncodeResponse,
	)
	
	http.Handle("/uppercase", uppercaseHandler)
	http.Handle("/count", countHandler)
	http.Handle("/metrics", promhttp.Handler())
	logger.Log("msg", "HTTP", "addr", ":8080")
	logger.Log("err", http.ListenAndServe(":8080", nil))
}
