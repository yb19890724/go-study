package main

import (
	"github.com/yb19890724/go-study/ddd/example1/task"
	"time"
)

func main() {

	task1 := new(task.Task1)
	go task1.Exec()
	task2 := new(task.Task2)
	go task2.Exec()
	task3 := new(task.Task3)
	go task3.Exec()
	task4 := new(task.Task4)
	go task4.Exec()
	time.Sleep(time.Second)
}
