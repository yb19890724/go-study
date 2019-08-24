package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main()  {
	
	// 随机因子 必须写 否则每次随机数都是一样的
	rand.Seed(time.Now().UnixNano())
	
	// 随机0，1 没有2 [0,n)
	// [0 大于等于
	// n) 小于
	fmt.Println(rand.Intn(2))
	fmt.Println(rand.Intn(2))
	fmt.Println(rand.Intn(2))
}
