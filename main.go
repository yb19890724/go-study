package main

import "fmt"

// 并发版求素数

func main() {
	origin, wait := make(chan int), make(chan struct{})
	Processor(origin,wait)
	
	for num :=2;num <10 ;num++{
		origin <-num
	}
	
	close(origin)
	<-wait
}

func Processor(seq chan int,wait chan struct{}){
	go func() {
		
		prime,ok :=<-seq
		
		if !ok{
			close(wait)
			return
		}
		fmt.Println(prime)
		
		out:=make(chan int )
		
		Processor(out ,wait)
		
		for num:=range seq{
			if num % prime!=0 {
				out <-num
			}
		}
		close(out)
		
	}()
}

