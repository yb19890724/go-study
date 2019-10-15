package configs

import (
	"context"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
	"github.com/yb19890724/go-study/consul/example3/pkg/consul"
	"log"
	"sync"
	"time"
)

var configs *configMap

type configMap struct {
	sync.RWMutex
	Conf map[string]DbConf
}

// mysql配置 结构体
type MysqlClusterConfig struct {
	Mysql MysqlConfig `json:"mysql"`
}

// 集群标识
type MysqlConfig struct {
	Cluster map[string]DbConf `json:"cluster"`
}

// db配置项
type DbConf struct {
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime"`
	Dsn             string        `json:"dsn"`
	MaxIdleConns    int           `json:"max_idle_conns"`
	MaxOpenConns    int           `json:"max_open_conns"`
}

func init() {
	configs = new(configMap)
	configs.Conf = make(map[string]DbConf)
}

// 获取指定配置
func (c *configMap) getConf(dbName string) bool {
	configs.RLock()
	_, ok := c.Conf[dbName]
	configs.RUnlock()
	return ok
}

// 设置配置
func (c *configMap) setConf(dbName string, conf DbConf) {
	configs.RLock()
	configs.Conf[dbName] = conf
	configs.RUnlock()
}

// reset configs all
func (c *configMap) resetConf(dbConfs map[string]DbConf) {
	
	for index,v :=range dbConfs {
		c.setConf(index,v)
	}
	fmt.Println("重置完毕")
}

// cn 配置名称
func GetConfig(dbName string, dbConfName string) (DbConf, error) {
	
	ok := configs.getConf(dbName)
	
	// 获取失败 从consul获取
	if !ok {
		
		cs := consul.NewConsulServer()
		
		// 获取配置
		_, conf, err := cs.Get(dbConfName)
		
		mConf, err := MysqlFormatConf(conf)
		
		// 读取失败
		if err != nil {
			return DbConf{}, err
		}
		
		configs.setConf(dbName, mConf.Mysql.Cluster[dbName])
		
	}
	
	return configs.Conf[dbName], nil
}

// 配置格式化
func MysqlFormatConf(conf []byte) (MysqlClusterConfig, error) {
	var mysqlConf MysqlClusterConfig
	
	err := yaml.Unmarshal(conf, &mysqlConf)
	
	if err != nil {
		log.Fatal(err)
	}
	return mysqlConf, err
}


// 监控配置 变化
func Watch(ctx context.Context,cancel context.CancelFunc,dbConfName string,changeConf chan<- int) {
	
	watchConfig := make(map[string]interface{})
	
	watchConfig["type"] = "key"
	watchConfig["key"] = dbConfName
	watchPlan, err := watch.Parse(watchConfig)
	
	if err != nil {
		fmt.Println(err)
	}
	
	watchPlan.Handler = func(idx uint64, data interface{}) {
		
		fmt.Println("配置发生变换")
		
		d,ok:=data.(*api.KVPair)
		
		if ok {
			mConf,err:=MysqlFormatConf(d.Value)
			// 重置配置
			if err == nil {
				fmt.Println("重置配置")
				configs.resetConf(mConf.Mysql.Cluster) // 重置配置
				changeConf<- 1
			}
		}
		
	}
	
	if err := watchPlan.Run("http://localhost:8500"); err != nil {
		log.Fatalf("start watch error, error message: %s", err.Error())
	}
	
	
	for {
		select {
		case <-ctx.Done():// 检测接受端是否关闭
			goto CLOSED
		}
	}
	
CLOSED:
	fmt.Println("监听配置变化 退出")
	// 关闭监听
	watchPlan.Stop()
	cancel()// 通知接收端关闭
}
