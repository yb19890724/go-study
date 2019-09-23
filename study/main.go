package main

import (
	"time"
)

const splitValue = 4294967296

// 分表
type Strategy struct {
	SplitType string // 分表类型
	RangeStartIndex int // 分表开始索引
	RangeEndIndex int // 分表开始索引
	MaxTableIndex int // 表最大索引值
	Reverse bool // 表最大索引值
	Limit bool // 表最大索引值
}


func timeById(value int) int {
	
	if value < splitValue {
		return 0
	}
	
	formatTime, _ := time.Parse("2006-01-02 15:04:05", "2015-01-01 00:00:00")
	
	time := int64((value>>21)/1000) + formatTime.Unix()
	
	return int(time)
}

func indexByTime(t int64) int64 {
	
	bt := "2015-11-11 00:00:00"
	invl := int64(864000)
	si := int64(24)
	moi := int64(24)
	
	formatTime, _ := time.Parse("2006-01-02 15:04:05", bt)
	
	timeGap := t - formatTime.Unix()
	
	if timeGap < 0 {
		return moi
	}
	return int64(timeGap)/invl + si
}

func (s *Strategy) buildRange(start int64, end int64, r bool, l bool) {
	
	s.Limit, s.Reverse = l, r
	
	mii := int64(0)
	cmi := time.Now().Unix()
	mxi := cmi
	
	if start > 0 {
		mii = indexByTime(start)
	}
	
	if end > 0 {
		mxi = indexByTime(end)
	}
	
	s.rangeStartIndex( mii, mxi)
	s.rangeEndIndex( mii, mxi)
	s.MaxTableIndex=int(cmi)
}

func (s *Strategy) rangeStartIndex( mi int64, mx int64)  {
	if s.Reverse == true {
		s.RangeStartIndex = int(mi)
	}
	s.RangeStartIndex = int(mx)
}

func (s *Strategy) rangeEndIndex( mi int64, mx int64)  {
	if s.Reverse == true {
		s.RangeEndIndex = int(mx)
	}
	s.RangeEndIndex = int(mi)
}

func main() {

}
