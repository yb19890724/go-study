package main

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	// 设置参数
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

// Fields
// 1)作用

func main() {
	event := "event"
	topic := "topic"
	key := 123123

	// 结构化信息记录，传统的记录方式如：
	log.Fatalf("Failed to send event %s to topic %s with key %d", event, topic, key)

	// 在logrus中不提倡这样写，鼓励使用Fields结构化日志内容，如：

	log.WithFields(log.Fields{
		"event": event,
		"topic": topic,
		"key":   key,
	}).Fatal("Failed to send event")

	// {"level":"fatal","msg":"Failed to send event event to topic topic with key 123123","time":"2019-07-12T16:28:21+08:00"}

}
