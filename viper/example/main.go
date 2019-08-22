package main

import (
	"fmt"
	"github.com/spf13/viper"
)

var (
	configMysql = "database"
	configRedis   = "redis"
)

// 定义config结构体
type redisConfig struct {
	Addr            string`json:"addr"`
	DialTimeout     string`json:"dialTimeout"`
	MaxConnAge      string`json:"maxConnAge"`
	MaxRetries      string`json:"maxRetries"`
	MinIdleConns    string`json:"minIdleConns"`
	PoolSize        string`json:"poolSize"`
	ReadTimeout     string`json:"readTimeout"`
	WriteTimeout    string`json:"writeTimeout"`
}

func main() {
	
	/*mysqlV := viper.New()
	mysqlV.SetConfigName(configMysql)
	mysqlV.AddConfigPath("$GOPATH/src/github.com/yb19890724/go-study/viper/example/config")
	
	err := mysqlV.ReadInConfig()
	if err != nil {
		fmt.Println(fmt.Errorf("Fatal error when reading %s config file:%s", configMysql, err))
	}
	
	dsn := mysqlV.GetString("AvatarDetail.master.dsn")
	fmt.Println(dsn)*/
	
	//---demo 2
	redisV := viper.New()
	redisV.SetConfigName(configRedis)
	redisV.AddConfigPath("$GOPATH/src/github.com/yb19890724/go-study/viper/example/config")
	
	err := redisV.ReadInConfig()
	if err != nil {
		fmt.Println(fmt.Errorf("Fatal error when reading %s config file:%s", configMysql, err))
	}
	
	addr := redisV.GetString("addr")
	fmt.Println(addr)
	var configJson redisConfig
	
	// 转换json
	
	if err := redisV.Unmarshal(&configJson); err != nil {
		fmt.Println(err)
	}
	fmt.Println(configJson.Addr)

}
