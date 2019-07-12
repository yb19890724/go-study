package main


// 逐行读取

import (
	"bufio"
	"io"
	"os"
)

func processLine(line []byte) {
	os.Stdout.Write(line) // 输出控制台
}

func ReadLine(filePth string, hookfn func([]byte)) error {
	f, err := os.Open(filePth)
	if err != nil {
		return err
	}
	defer f.Close()
	
	bfRd := bufio.NewReader(f)
	for {
		line, err := bfRd.ReadBytes('\n') // 根据换行标记切分
		hookfn(line) // 放在错误处理前面，即使发生错误，也会处理已经读取到的数据。
		if  err != nil { // 遇到任何错误立即返回，并忽略 EOF 错误信息
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
	
	return nil
}




func main() {
	ReadLine("./read/test.txt", processLine)
}