package main

import (
	"bytes"
	"fmt"
	
)

// go 语言里的合集一般会用map[t]bool 这种形式来表示，T代表元素类型。集合用map类型来标识虽然非常灵活，但我们可以一种更好的形式来标识它，
// 例如在数据流分析领域，集合元素通常是一个非负整数，集合会包含很多元素，并且集合经常会进行并集，交集操作，这种情况下，bit数组会比map表现更加理想。
// 解释：比如我们执行一个http下载任务，把文件按照16kb一块划分为很多块，需要有一个全局变量来标识那些块下载完成了，这种时候也需要bit数组
// 一个bit数组通常会用一个无符号或者称之为slice，每个元素的每一位都表示集合里的一个值，当集合的第i位被设置时，我们才说这个集合包含元素i

type IntSet struct {
	words []uint64
}
// 报告该集合是否包含非负值x。
func (s *IntSet) Has(x int) bool {
	word, bit := x/64, uint(x%64)
	
	fmt.Println(s.words[word])
	fmt.Println((1<<bit))
	fmt.Println(s.words[word]&(1<<bit))
	return word < len(s.words) && s.words[word]&(1<<bit) != 0
}
// 添加将非负值x添加到集合中。
func (s *IntSet) Add(x int) {
	word, bit := x/64, uint(x%64)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
}
// UnionWith将s设置为s和t的并集。
func (s *IntSet) UnionWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] |= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

// String将集合返回为“{1 2 3}”形式的字符串。
func (s *IntSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < 64; j++ {
			if word&(1<<uint(j)) != 0 {
				if buf.Len() > len("{") {
					buf.WriteByte('}')
				}
				fmt.Fprintf(&buf, "%d", 64*i+j)
			}
		}
	}
	buf.WriteByte('}')
	return buf.String()
}

func main() {
	var x, _ IntSet
	x.Add(1)
	x.Add(144)
	x.Add(9)
	fmt.Println(x.Has(9)) // "{1 9 144}"
	
	
	// y.Add(9)
	// y.Add(42)
	// fmt.Println(y.String()) // "{9 42}"
	// x.UnionWith(&y)
	// fmt.Println(x.String()) // "{1 9 42 144}"
	//
	//
	// fmt.Println(x.Has(9), x.Has(123)) // "true false"
}
