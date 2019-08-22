package main

import (
	"github.com/go-kit/kit/log"
	"os"
)

func main() {
	
	// demo 1
	var logger log.Logger
	{
		logger =log.NewLogfmtLogger(os.Stdout)
		logger =log.WithPrefix(logger,"my test","1.0")
		logger =log.With(logger,"time",log.DefaultTimestampUTC)
		logger =log.With(logger,"caller",log.DefaultCaller)
	}
	
	logger.Log("method","get")
	
	
}
