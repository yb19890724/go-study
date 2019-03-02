package main

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"math/rand"
	"time"
)

type lotteryController struct {
	Ctx iris.Context
}

// 创建应用
func newApp() *iris.Application {
	
	app := iris.New() // 加载应用
	
	mvc.New(app.Party("/")).Handle(&lotteryController{})
	
	return app
}

func main() {
	
	app := newApp()
	
	app.Run(iris.Addr(":8080")) // 启动服务
	
}

// 即开即得型 http://localhost:8000
func (c *lotteryController) Get() string {
	
	var prize string
	
	seed := time.Now().UnixNano()
	code := rand.New(rand.NewSource(seed)).Intn(10)
	
	switch {
	
	case 1 == code:
		prize = "一等奖"
	
	case 2 == code:
		prize = "二等奖"
	
	case 3 == code:
		prize = "三等奖"
	
	default:
		return fmt.Sprintf("很遗憾您未中奖，兑换码：%s 谢谢参与!", code)
	}
	
	return fmt.Sprintf("尾号为1获得一等奖\n"+
		"尾号为2获得二等奖\n"+
		"尾号为3获得三等奖\n"+
		"兑换码：%s,恭喜您获得%s", code, prize)
}

// 双色球自选型
func (c *lotteryController) GetPrize() string {
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	
	var prize [7]int
	
	// 6个红色球 1-33
	for i := 0; i < 6; i++ {
		prize[i] = r.Intn(33) + 1
	}
	
	// 最后一位蓝球 1-16
	prize[6] = r.Intn(16) + 1
	
	return fmt.Sprintf("今日开奖号码是：%v", prize)
}
