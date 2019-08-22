package main

import (
	"flag"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/config"
	"github.com/prometheus/common/log"
	"os"
)

type Config struct {
	Version string
	Hello   struct {
		Name string
	}
	Etcd struct {
		Addrs    []string
		UserName string
		Password string
	}
}

func main() {

	testFlag := cli.StringFlag{
		Name:  "f",
		Value: os.Getenv("GOPATH") + "/src/github.com/yb19890724/go-study/micro/example4/config/config.json",
		Usage: "please use config.json",
	}

	//configFlag := flag.String("f", os.Getenv("GOPATH")+"/src/github.com/yb19890724/go-study/micro/example4/config/config.json", "please use config.json")
	configFlag := flag.String(testFlag.Name, testFlag.Value, testFlag.Usage)
	conf := new(Config)

	if err := config.LoadFile(*configFlag); err != nil {
		log.Fatal(err)

	}

	if err := config.Scan(conf); err != nil {
		log.Fatal(err)
	}

	micro.NewService(
		micro.Flags(testFlag),
	)

	println(conf.Version)
}
