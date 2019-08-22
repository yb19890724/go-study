package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

// 一次性获取文件所有内容

// ioutil.ReadFile(filePth) 更简单方式
func ReadAll(filePth string) ([]byte, error) {
	f, err := os.Open(filePth)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(f)
}

// 再看看输出：可以看到len只有9但是cap却有1536，
// 我们只读取很少的内容却使用了这么多的内存，这在平时不会有问题，
// 但是比如在网络应用当有大量请求过来时就容易导致内存严重浪费，严重时还会内存泄漏。

//溯源
//我们来看看它底层到底如何读取的：ReadAll调用了内部方法readAll

func main() {
	fileInfo, _ := ReadAll("./read/test.txt")
	fmt.Println(string(fileInfo))
	fmt.Println(len(fileInfo)) // 9
	fmt.Println(cap(fileInfo)) // 16
}
