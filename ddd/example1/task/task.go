package task

import (
	"fmt"
	"github.com/yb19890724/go-study/ddd/example1/eventhandlers"
	"time"
)

const (
	task1HasCompleted eventhandlers.Event = "task1 has completed"
	task2HasCompleted eventhandlers.Event = "task2 has completed"
)

type Task1 struct {
}

type Task1Handler struct {
}

// 执行 触发处理方法
func (t *Task1) Exec() {
	task1Handler := new(Task1Handler)
	task1Handler.Handle()
}

func (t *Task1Handler) Handle() {
	fmt.Println("task1 handler start")
	time.Sleep(50 * time.Millisecond)
	fmt.Println("task1 handler end")
	ehs := eventhandlers.GetInstance()
	ehs.Pub(task1HasCompleted)// 发布通知已完成
	fmt.Println("task1 pub task1HasCompleted")
}

type Task2 struct {
}

type Task2Handler struct {
}

func (t *Task2) Exec() {
	task2Handler := new(Task2Handler)
	ehs := eventhandlers.GetInstance()
	ehs.Sub(task1HasCompleted, task2Handler)
	fmt.Println("task2 sub task1HasCompleted")
}

func (t *Task2Handler) Handle() {
	fmt.Println("task2 handler start")
	time.Sleep(100 * time.Millisecond)
	fmt.Println("task2 handler end")
	ehs := eventhandlers.GetInstance()
	ehs.Pub(task2HasCompleted)
	fmt.Println("task2 pub task2HasCompleted")
}

type Task3 struct {
}

type Task3Handler struct {
}

func (t *Task3) Exec() {
	task3Handler := new(Task3Handler)
	ehs := eventhandlers.GetInstance()
	ehs.Sub(task2HasCompleted, task3Handler)
	fmt.Println("task3 sub task2HasCompleted")
}

// 处理函数
func (t *Task3Handler) Handle() {
	fmt.Println("task3 handler start")
	time.Sleep(200 * time.Millisecond)
	fmt.Println("task3 handler end")
}

type Task4 struct {
}

type Task4Handler struct {
}

func (t *Task4) Exec() {
	task4Handler := new(Task4Handler)
	ehs := eventhandlers.GetInstance()
	ehs.Sub(task2HasCompleted, task4Handler)
	fmt.Println("task4 sub task2HasCompleted")
}

func (t *Task4Handler) Handle() {
	fmt.Println("task4 handler start")
	time.Sleep(200 * time.Millisecond)
	fmt.Println("task4 handler end")
}
