package main

import (
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"math/rand"
	"sync"
	"time"
)

type Product struct {
	ID    int
	Title string
	Price int
}

func getProduct() (Product, error) {
	r := rand.Intn(10)
	if r < 6 {
		time.Sleep(time.Second * 3)
	}
	return Product{
		ID:    101,
		Title: "Golang从入门到精通",
		Price: 12,
	}, nil
}

// 类似缓存数据源，专门来做服务降级用
// 一旦获取商品超时，返回一个推荐商品
func RecProduct() (Product, error) {
	return Product{
		ID :99,
		Title: "推荐商品",
		Price: 22,
	},nil
}

func main() {
	rand.Seed(time.Now().UnixNano())
	
	hystrix.ConfigureCommand("get_prod2", hystrix.CommandConfig{
		Timeout: 2000,
		MaxConcurrentRequests:5,// 最大 并发数
		RequestVolumeThreshold: 3,  // 请求阈值 默认20，(总请求值达到20个才进行错误百分比计算)
		ErrorPercentThreshold:50,// 错误百分比(默认50 50%) 达到多少的时候打开降级服务，不会执行原方法，直接执行降级服务
	})
	
	resultChan := make(chan Product,1)
	
	wg :=sync.WaitGroup{}
	
	for i:=0;i<10;i++{ // 修改i执行范围 发生改变
		go func() {
			wg.Add(1)
			defer wg.Done()
			// 执行过程 包装成一个command
			// name     : 名称
			// runFunc  : 运行逻辑
			// fallback : 降级服务
			
			// 开启goruntine 异步执行
			errs := hystrix.Go("get_prod2", func() error {
				// 里面代表业务逻辑
				p, err := getProduct() // 这里会随机延迟三秒
				
				if err != nil {
					return err
				}
				resultChan<-p
				return nil
			}, func(e error) error { // 服务降级
				ret, err := RecProduct()
				resultChan <- ret
				fmt.Println(e)
				return err
			})
			
			
			
			select {// 获取商品
			case getPrd:=<-resultChan:
				fmt.Println(getPrd)
			case err:=<-errs: // 获取错误
				fmt.Println(err)
			}
			
			
			time.Sleep(time.Second * 1)
			
		}()
	
	}
	wg.Wait()
}
