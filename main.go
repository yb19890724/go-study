package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT,syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGUSR1)
	
	for {
		s := <-ch
		switch s {
		case syscall.SIGINT:
			fmt.Println("SIGINT")
			return
		case syscall.SIGQUIT:
			fmt.Println("SIGQUIT")
			return
		case syscall.SIGSTOP:
			fmt.Println("SIGSTOP")
			return
		case syscall.SIGHUP:
			fmt.Println("SIGHUP")
			return
		case syscall.SIGKILL:
			fmt.Println("SIGKILL")
			return
		case syscall.SIGUSR1:
			fmt.Println("SIGUSR1")
			return
		default:
			fmt.Println("default")
			return
		}
	}
}

