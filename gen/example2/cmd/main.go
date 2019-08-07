package main

import (
	"github.com/yb19890724/go-study/gen/example2/pkg/option"
	"log"
	"os"
	
	"github.com/urfave/cli"
)

func main() {
	
	opt := option.Option{}
	
	app := cli.NewApp()
	
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "f",
			Value:       "./test.proto", // 默认值
			Usage:       "指定文件",
			Destination: &opt.Proto3Filename,
		},
		cli.StringFlag{
			Name:        "p",
			Usage:       "输出目录",
			Destination: &opt.Output,
		},
	}
	
	// 回调函数
	app.Action = func(c *cli.Context) error {
		return nil
	}
	
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
