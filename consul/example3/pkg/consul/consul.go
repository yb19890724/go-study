package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"log"
)

type ConsulServer struct {
	Conn *api.Client
}

func NewConsulServer() *ConsulServer {
	
	// register consul
	consulConfig := api.DefaultConfig()
	consulConfig.Address = "127.0.0.1:8500"
	client, err := api.NewClient(consulConfig)
	if err != nil {
		fmt.Printf("NewClient error\n%v", err)
		return nil
	}
	
	cs := &ConsulServer{ // 存储consul 服务
		Conn: client,
	}
	return cs
}



func (cs *ConsulServer) Put(k string, v []byte) {
	
	pair := &api.KVPair{Key: k, Value: v}
	if _, err := cs.Conn.KV().Put(pair, nil); err != nil {
		log.Fatal(err)
	} else {
		log.Println("put key ", k, "=>", v, " successfully")
	}
}

func (cs *ConsulServer) Get(k string) (string, []byte, error) {
	
	p, _, err := cs.Conn.KV().Get(k, nil)
	
	if err != nil {
		log.Fatal(err)
	}
	// log.Fatal(p.Key, "=>", p.Value)
	return p.Key, p.Value, err
}

func (cs *ConsulServer) List(k string) {
	if ps, _, err := cs.Conn.KV().List(k, nil); err != nil {
		log.Fatal(err)
	} else {
		for _, p := range ps {
			log.Println(p.Key, "=>", string(p.Value))
		}
	}
}

func (cs *ConsulServer) Keys(k string) {
	if ks, _, err := cs.Conn.KV().Keys(k, "/", nil); err != nil {
		log.Fatal(err)
	} else {
		for _, k := range ks {
			log.Println(k)
		}
	}
}

func (cs *ConsulServer) Delete(k string) {
	if _, err := cs.Conn.KV().Delete(k, nil); err != nil {
		log.Fatal(err)
	} else {
		log.Println("delete key '", k, "' successfully")
	}
}
