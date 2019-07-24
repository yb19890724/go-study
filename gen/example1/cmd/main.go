package main

import (
	
	"github.com/gogo/protobuf/proto"
	test "github.com/yb19890724/go-study/gen/example1/proto"
	"io/ioutil"
	"log"
)

func write() {
	
	
	c1 := &test.Class{
		Num: 1,
		Students: []*test.Student{
			{Name: "xiaoming", Age: 21, Sex: test.Sex_MAN},
			{Name: "xiaohua", Age: 21, Sex: test.Sex_WOMAN},
			{Name: "xiaojin", Age: 21, Sex: test.Sex_MAN},
		},
	
	}
	
	c1.Reset()
	
	// 使用protobuf工具把struct数据类型格式化成字节数组（压缩和编码）
	data, _ := proto.Marshal(c1)
	log.Println(string(data))
	// 把字节数组写入到文件中
	//ioutil.WriteFile("test.txt", data, os.ModePerm)
}

func read() {
	// 以字节数组的形式读取文件内容
	data, _ := ioutil.ReadFile("test.txt")
	
	class := new(test.Class)
	
	// 使用protobuf工具把字节数组解码成struct(解码)
	proto.Unmarshal(data, class)
	
	log.Println(class.Num)
	for _, v := range class.Students {
		log.Println(v.Name, v.Age, v.Sex)
	}
}

func main() {
	write()
	read()
}