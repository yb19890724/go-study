package main

import (
	"context"
	"flag"
	`github.com/yb19890724/go-study/gokit/pkg/stringsvc3`
	"net/http"
	"os"

	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	httptransport "github.com/go-kit/kit/transport/http"
)

func main() {
	var (
		listen = flag.String("listen", ":8082", "HTTP listen address")
		proxy  = flag.String("proxy", "localhost:8080,localhost:8081", "Optional comma-separated list of URLs to proxy uppercase requests")
	)
	flag.Parse()

	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "listen", *listen, "caller", log.DefaultCaller)

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
	}, []string{})

	var svc stringsvc3.StringService
	svc = stringsvc3.Service{}
	svc = stringsvc3.ProxyingMiddleware(context.Background(), *proxy, logger)(svc)
	svc = stringsvc3.LoggingMiddleware(logger)(svc)
	svc = stringsvc3.InstrumentingMiddleware(requestCount, requestLatency, countResult)(svc)

	uppercaseHandler := httptransport.NewServer(
		stringsvc3.MakeUppercaseEndpoint(svc),
		stringsvc3.DecodeUppercaseRequest,
		stringsvc3.EncodeResponse,
	)
	
	countHandler := httptransport.NewServer(
		stringsvc3.MakeCountEndpoint(svc),
		stringsvc3.DecodeCountRequest,
		stringsvc3.EncodeResponse,
	)

	http.Handle("/uppercase", uppercaseHandler)
	http.Handle("/count", countHandler)
	http.Handle("/metrics", promhttp.Handler())
	logger.Log("msg", "HTTP", "addr", *listen)
	logger.Log("err", http.ListenAndServe(*listen, nil))
}
