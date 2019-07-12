package main

import (
	"bufio"
	"fmt"
	"os"
)

// 并发 读取文件

func main() {
	
	
	f, err :=os.OpenFile("./read/test.txt",os.O_APPEND,0644)
	
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	
	
	bfRd := bufio.NewReader(f)
	
	ret,err:=f.Seek(1,os.SEEK_SET)
	fmt.Println(ret)
	line, err := bfRd.ReadBytes('\n') // 根据换行标记切分
	fmt.Println(string(line))
	
	
	
/*
	bfRd := bufio.NewReader(f)

	line, err := bfRd.ReadBytes('\n') // 根据换行标记切分
	fmt.Println(string(line))
	*/
}