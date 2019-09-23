package main

import (
	"fmt"
	"log"

	"net/http"

	consulapi "github.com/hashicorp/consul/api"
)

func consulCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "consulCheck")
}

func registerServer() {

	config := consulapi.DefaultConfig()
	client, err := consulapi.NewClient(config)

	if err != nil {
		log.Fatal("consul client error : ", err)
	}

	checkPort := 8080

	registration := new(consulapi.AgentServiceRegistration)
	registration.ID = "serverNode_1"// 服务节点的名称
	registration.Name = "serverNode"// 服务名称
	registration.Port = 9527 // 服务端口
	registration.Tags = []string{"serverNode"}// tag，可以为空
	registration.Address = "127.0.0.1"
	registration.Check = &consulapi.AgentServiceCheck{ // 健康检查
		HTTP:                           fmt.Sprintf("http://%s:%d%s", registration.Address, checkPort, "/check"),
		Timeout:                        "3s",// 超时时间
		Interval:                       "5s",// 健康检查间隔
		DeregisterCriticalServiceAfter: "30s",// check失败后30秒删除本服务  注销时间，相当于过期时间
	}

	err = client.Agent().ServiceRegister(registration) // 注册服务

	if err != nil {
		log.Fatal("register server error : ", err)
	}

	http.HandleFunc("/check", consulCheck)
	http.ListenAndServe(fmt.Sprintf(":%d", checkPort), nil)

}

func main() {
	registerServer()
}
