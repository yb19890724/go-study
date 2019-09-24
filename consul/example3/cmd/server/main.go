package main

import (
	"fmt"
	"github.com/yb19890724/go-study/consul/example3/pkg/storage"
	"time"
)

type User struct {
	ID   int    `json:"id"`   // 列名为 `id`
	Name string `json:"name"` // 列名为 `username`
}

var Users []User

func main() {
	
	// watch k/v数据变换更新连接
	go func() {
		storage.Watch("/cluster/database")
	}()
	
	
	// 注意这是个例子，延迟2秒是因为 重置连接会关闭 数据库连接，如果是1秒直接就退出程序了
	// 正确方式应该是一个请求一个连接，更新配置了删除旧连接，关闭连接。
	// 再次到达的请求使用新的连接来连接数据库
	for {
		
		time.Sleep(1*time.Second)
		
		db, _ := storage.GetMysqlConnection("master", "/cluster/database")
		
		if err:=db.Select([]string{"id", "name"}).First(&Users, 2).Error; err != nil {
			fmt.Println(err)
			return
		}
		
		fmt.Println(Users)
	}
	
}
