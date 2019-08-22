package main

import (
	"fmt"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/yb19890724/go-study/gokit/pkg/stringsvc1"
	"net/http"
	"os"
)

func main() {

	logger := log.NewLogfmtLogger(os.Stderr)
	mid := stringsvc1.LoggingMiddleware(log.With(logger, "LOGMETHOD", "uppercase"))

	svc := stringsvc1.Service{}

	uppercaseHandler := httptransport.NewServer(
		mid(stringsvc1.MakeUppercaseEndpoint(svc)),
		stringsvc1.DecodeUppercaseRequest,
		stringsvc1.EncodeResponse,
	)

	countHandler := httptransport.NewServer(
		mid(stringsvc1.MakeCountEndpoint(svc)),
		stringsvc1.DecodeCountRequest,
		stringsvc1.EncodeResponse,
	)

	http.Handle("/uppercase", uppercaseHandler)
	http.Handle("/count", countHandler)
	http.ListenAndServe(":8080", nil)
	fmt.Println("localhost:8080")
}
