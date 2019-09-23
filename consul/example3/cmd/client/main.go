package main

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
	"github.com/yb19890724/go-study/consul/example3/pkg/storage"
	"log"
)

func main() {
	
	newWatch()
}

func newWatch() {
	
	watchConfig := make(map[string]interface{})
	
	watchConfig["type"] = "key"
	watchConfig["key"] = "/cluster/database"
	watchPlan, err := watch.Parse(watchConfig)
	
	if err != nil {
		fmt.Println(err)
	}
	
	watchPlan.Handler = func(idx uint64, data interface{}) {
		
		d:=data.(*api.KVPair)
		
		mConf,err := storage.FormatConf(d.Value)
		
		fmt.Println(err)
		fmt.Println(mConf)
		
	}
	
	if err := watchPlan.Run("http://localhost:8500"); err != nil {
		log.Fatalf("start watch error, error message: %s", err.Error())
	}
	
	watchPlan.Stop()
}
