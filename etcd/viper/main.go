package main

import (
	"fmt"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

func main() {
	
	
	viper.AddRemoteProvider("etcd", "localhost:2379","/configs/hugo.json")
	viper.SetConfigType("json") // because there is no file extension in a stream of bytes, supported extensions are "json", "toml", "yaml", "yml", "properties", "props", "prop"
	err := viper.ReadRemoteConfig()
	
	fmt.Println(err)
}
