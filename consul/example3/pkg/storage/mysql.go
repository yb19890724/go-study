package storage

import (
	"context"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/yb19890724/go-study/consul/example3/pkg/configs"
	"log"
	"sync"
	"time"
)

var dbs *Dbs


type Dbs struct {
	sync.RWMutex
	Map map[string]*gorm.DB
}


func init() {
	dbs = new(Dbs)
	dbs.Map = make(map[string]*gorm.DB)
}

// 获取指定 db连接
func (d *Dbs) getDbsConn(dbName string) bool {
	dbs.RLock()
	_, ok := dbs.Map[dbName]
	dbs.RUnlock()
	return ok
}

// 设置db连接
func (d *Dbs) setDbsConn(dbName string, db *gorm.DB) {
	dbs.Lock()
	dbs.Map[dbName] = db
	dbs.Unlock()
}

// 设置db连接
func (d *Dbs) delDbsConn(dbName string) {
	dbs.Lock()
	
	delete(d.Map,dbName)
	dbs.Unlock()
}


// c 读取对应配置名称
func GetMysqlConnection(dbName string, dbConfName string) (*gorm.DB, error) {
	
	ok := dbs.getDbsConn(dbName)
	if !ok {
		
		// 获取配置
		conf, err := configs.GetConfig(dbName, dbConfName)
		
		if err != nil {
			log.Fatal(err)
		}
		
		// 创建连接
		dbConn, err := createMysqlConnection(conf)
		
		if err != nil {
			return nil, err
		}
		
		dbs.setDbsConn(dbName, dbConn)
		
	}
	return dbs.Map[dbName], nil
}

// 监听k/v变化
func Watch (dbConfName string)  {
	
	changConf:=make(chan int,1)
	
	ctx, cancel := context.WithCancel(context.Background())
	
	
	go func() {
		dbs.resetDbsConn(ctx, cancel, changConf)
	}()
	
	configs.Watch(ctx, cancel, dbConfName, changConf)

}

// 卸载db连接
func (d *Dbs) resetDbsConn(ctx context.Context,cancel context.CancelFunc,changConf <-chan int)  {
	
	for {
		select {
		case <-changConf:
			fmt.Println("重置db")
			
			for index, v := range d.Map{
				
				d.delDbsConn(index)
				
				// 这里等待是因为，如果请求正在使用db连接直接关闭会报错，这时等待调用结束后在进行关闭
				// 可以默认判断5秒钟请求如果没有结束，证明是一个不好的请求，所以直接关闭就行了、
				time.Sleep(5*time.Second)
				
				err :=v.Close()// 关闭数据库连接
				
				if err !=nil {
					log.Printf("db close err %s",err)
					goto ERROR
				}
			}
		case <-ctx.Done():
			goto CLOSED
		}
	}
ERROR:
	cancel()
CLOSED:
	fmt.Println("重置机制退出")
}


// 创建mysql连接
func createMysqlConnection(conf configs.DbConf) (*gorm.DB, error) {
	
	db, err := gorm.Open("mysql", conf.Dsn)
	
	if err != nil {
		return nil, errors.New(fmt.Sprintf("connect mysql fail %s", err))
	}
	// 设置连接池
	db.DB().SetMaxIdleConns(conf.MaxIdleConns)
	db.DB().SetMaxOpenConns(conf.MaxOpenConns)
	db.DB().SetConnMaxLifetime(conf.ConnMaxLifetime)
	db.SingularTable(true)
	db.BlockGlobalUpdate(false)
	return db, nil
}
