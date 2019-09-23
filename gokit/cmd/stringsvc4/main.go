package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"github.com/nats-io/nats.go"
	"log"
	"net/http"
	"strings"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	natstransport "github.com/go-kit/kit/transport/nats"

)

// StringService provides operations on strings.
type StringService interface {
	Uppercase(context.Context, string) (string, error)
	Count(context.Context, string) int
}

// stringService is a concrete implementation of StringService
type stringService struct{}

func (stringService) Uppercase(_ context.Context, s string) (string, error) {
	if s == "" {
		return "", ErrEmpty
	}
	return strings.ToUpper(s), nil
}

func (stringService) Count(_ context.Context, s string) int {
	return len(s)
}

// ErrEmpty is returned when an input string is empty.
var ErrEmpty = errors.New("empty string")

// For each method, we define request and response structs
type uppercaseRequest struct {
	S string `json:"s"`
}

type uppercaseResponse struct {
	V   string `json:"v"`
	Err string `json:"err,omitempty"` // errors don't define JSON marshaling
}

type countRequest struct {
	S string `json:"s"`
}

type countResponse struct {
	V int `json:"v"`
}

// 端点是Go工具包中的主要抽象。端点表示单个RPC（服务接口中的方法）
func makeUppercaseHTTPEndpoint(nc *nats.Conn) endpoint.Endpoint {
	return natstransport.NewPublisher(
		nc,
		"stringsvc.uppercase",
		natstransport.EncodeJSONRequest,
		decodeUppercaseResponse,
	).Endpoint()
}

func makeCountHTTPEndpoint(nc *nats.Conn) endpoint.Endpoint {
	return natstransport.NewPublisher(
		nc,
		"stringsvc.count",
		natstransport.EncodeJSONRequest,
		decodeCountResponse,
	).Endpoint()
}

// http 服务
func makeUppercaseEndpoint(svc StringService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(uppercaseRequest)
		v, err := svc.Uppercase(ctx, req.S)
		if err != nil {
			return uppercaseResponse{v, err.Error()}, nil
		}
		return uppercaseResponse{v, ""}, nil
	}
}

func makeCountEndpoint(svc StringService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(countRequest)
		v := svc.Count(ctx, req.S)
		return countResponse{v}, nil
	}
}

// 传输向网络公开服务。在第四个示例中，我们使用了基于NATS和HTTP的JSON。
func main() {
	svc := stringService{}

	natsURL := flag.String("nats-url", nats.DefaultURL, "URL for connection to NATS")
	flag.Parse()

	nc, err := nats.Connect(*natsURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	uppercaseHTTPHandler := httptransport.NewServer(
		makeUppercaseHTTPEndpoint(nc),
		decodeUppercaseHTTPRequest,
		httptransport.EncodeJSONResponse,
	)

	countHTTPHandler := httptransport.NewServer(
		makeCountHTTPEndpoint(nc),
		decodeCountHTTPRequest,
		httptransport.EncodeJSONResponse,
	)

	uppercaseHandler := natstransport.NewSubscriber(
		makeUppercaseEndpoint(svc),
		decodeUppercaseRequest,
		natstransport.EncodeJSONResponse,
	)

	countHandler := natstransport.NewSubscriber(
		makeCountEndpoint(svc),
		decodeCountRequest,
		natstransport.EncodeJSONResponse,
	)

	uSub, err := nc.QueueSubscribe("stringsvc.uppercase", "stringsvc", uppercaseHandler.ServeMsg(nc))
	if err != nil {
		log.Fatal(err)
	}
	defer uSub.Unsubscribe()

	cSub, err := nc.QueueSubscribe("stringsvc.count", "stringsvc", countHandler.ServeMsg(nc))
	if err != nil {
		log.Fatal(err)
	}
	defer cSub.Unsubscribe()

	http.Handle("/uppercase", uppercaseHTTPHandler)
	http.Handle("/count", countHTTPHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func decodeUppercaseHTTPRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request uppercaseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeCountHTTPRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request countRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeUppercaseResponse(_ context.Context, msg *nats.Msg) (interface{}, error) {
	var response uppercaseResponse

	if err := json.Unmarshal(msg.Data, &response); err != nil {
		return nil, err
	}

	return response, nil
}

func decodeCountResponse(_ context.Context, msg *nats.Msg) (interface{}, error) {
	var response countResponse

	if err := json.Unmarshal(msg.Data, &response); err != nil {
		return nil, err
	}

	return response, nil
}

func decodeUppercaseRequest(_ context.Context, msg *nats.Msg) (interface{}, error) {
	var request uppercaseRequest

	if err := json.Unmarshal(msg.Data, &request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeCountRequest(_ context.Context, msg *nats.Msg) (interface{}, error) {
	var request countRequest

	if err := json.Unmarshal(msg.Data, &request); err != nil {
		return nil, err
	}
	return request, nil
}
