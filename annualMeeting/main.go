package main

/**
 * curl --data "users=test1,test2" http://localhost:8080/import
 * curl http://localhost:8080
 * curl http://localhost:8080/lucky
 */
import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"math/rand"
	"strings"
	"sync"
	"time"
)

var (
	userList []string
	mutex    sync.Mutex // 互斥锁
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

	userList = []string{}

	// 互斥锁初始化
	mutex = sync.Mutex{}

	app.Run(iris.Addr(":8080")) // 启动服务

}

// 获取当前参与抽奖人数
// @method post : http://localhost:8080
// @return string
func (c *lotteryController) Get() string {

	count := len(userList)

	return fmt.Sprintf("当前总共参与抽奖的用户数：%d\n", count)

}

// 导入抽奖人员
// @method : post http://localhost:8080/import
// @params : users
func (c *lotteryController) PostImport() string {

	strUsers := c.Ctx.FormValue("users")

	users := strings.Split(strUsers, ",")

	// 上锁
	mutex.Lock()

	defer mutex.Unlock()

	beforeCount := len(userList)
	for _, u := range users { // 循环导入用户
		u = strings.TrimSpace(u)
		if len(u) > 0 { // 判断空不导入
			userList = append(userList, u)
		}
	}
	afterCount := len(userList) // 导入后用户数量

	return fmt.Sprintf("当前总共参与抽奖的用户数：%d，成功导入的用户数：%d\n",
		afterCount, (afterCount - beforeCount))

}

// 抽奖方法
// @method get : http://localhost:8080/lucky
func (c *lotteryController) GetLucky() string {

	// 上锁
	mutex.Lock()

	defer mutex.Unlock()

	count := len(userList)

	// 抽奖人数大于1
	if count > 1 {

		// 时间戳
		seed := time.Now().UnixNano()

		// 随机数，在count范围内
		index := rand.New(rand.NewSource(seed)).Int31n(int32(count))

		user := userList[index]

		// 取出当前随机数index之前所有和当前index后面所有
		userList = append(userList[0:index], userList[index+1:]...)

		return fmt.Sprintf("当前中奖用户: %s,剩余抽奖人数%d\n", user, count-1)

	} else if count == 1 { // 抽奖人数等于1

		user := userList[0:0]

		return fmt.Sprintf("当前中奖用户: %s,剩余抽奖人数%d\n", user, count-1)

	}

	return fmt.Sprintf("已经没有参与用户，请先添加抽奖人员")
}
