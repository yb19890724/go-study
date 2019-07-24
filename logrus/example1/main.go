package main


// logrus有一个默认的Logger
// 当一个应用中，需要向多个地方输出时，需要不同的Logger实例
import (
	"github.com/sirupsen/logrus"
	"os"
)

// 可以使用以上语句，创建Logger实例
var log = logrus.New()

func main() {
	// 创建文件
	file ,err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil{
		log.Out = file
	}else{
		log.Info("Failed to log to file")
	}
	
	// 写入文件内容
	log.WithFields(logrus.Fields{
		"filename": "123.txt",
	}).Info("打开文件失败")
}

// 根目录./logrus.log 文件 写入 日志内容