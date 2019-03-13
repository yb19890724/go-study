package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"
)

type LogProcess struct {
	readChannel  chan []byte //读取管道
	writeChannel chan string //写入管道
	read         Reader      //定义接收读接口
	write        Writer      //定义接收写接口
}

type Reader interface {
	Read(readChannel chan []byte)
}

type Writer interface {
	Write(writeChannel chan string)
}

type ReadFile struct {
	path string // 读取路径
}

type WriteFile struct {
	path string
}

// 读取模块
func (r *ReadFile) Read(readChannel chan []byte) {

	f, err := os.Open(r.path)

	if nil != err { //打开文件是否正常
		panic(fmt.Sprintf("oepn file errpr %s", err.Error()))
	}

	// 从文件末尾开始逐行读取文件内容 (读取最新的内容)
	f.Seek(0, 2)
	read := bufio.NewReader(f)

	for {

		line, err := read.ReadBytes('\n')

		if err == io.EOF {
			time.Sleep(500 * time.Millisecond)
			continue
		} else if err != nil {
			panic(fmt.Sprintf("ReadBytes error:%s", err.Error()))
		}

		readChannel <- line[:len(line)-1]
	}

}

//解析模块
func (l *LogProcess) Process() {

	for v := range l.readChannel {
		l.writeChannel <- string(v)
	}

}

//写入模块
//用引用方式不用做拷贝，速度上有优势
//使用l修改自身的参数
func (w *WriteFile) Write(writeChannel chan string) {
	//循环读取

	for v := range writeChannel {
		fmt.Println(fmt.Sprintf("%s", v))
	}
}

func main() {

	fmt.Print(1)


	read := &ReadFile{
		path: "/Users/huangyibing/Home/access.log",
	}

	write := &WriteFile{
		path: "D:\\phpStudy\\PHPTutorial\\nginx\\logs\access.log",
	}


	lp := &LogProcess{
		readChannel:  make(chan []byte),
		writeChannel: make(chan string),
		read:         read,
		write:        write,
	}

	go lp.read.Read(lp.readChannel)
	go lp.Process()
	go lp.write.Write(lp.writeChannel)

	time.Sleep(30 * time.Second)
}
