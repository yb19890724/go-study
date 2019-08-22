package main

import (
	"fmt"
	"log"
	"time"
)

// 限流桶
// 限流间隔时间
// 速率
// 容量

var (
	quantum         = int64(1)
	rate            = float64(2)
	capacity        = int64(2)
	count           = int64(1)
	availableTokens = int64(2)
	startTime       = time.Now()
)

func bucket() {

	var lastetTick int64
	t := time.NewTicker(500 * time.Millisecond)
	fillInterval := time.Duration(1e9 * float64(quantum) / rate)

	for e := range t.C {

		tick := int64(time.Now().Sub(startTime) / fillInterval)

		if availableTokens >= capacity {

			fmt.Printf("ok use [%d]", availableTokens)

		} else {

			availableTokens += (tick - lastetTick) * quantum

			if availableTokens > capacity {
				availableTokens = capacity
			}

			lastetTick = tick
		}

		if availableTokens <= 0 {

			log.Println("no availableTokens")

		} else {

			if count > availableTokens {

				count = availableTokens

			}
			availableTokens -= count

			fmt.Printf("[%+v] === >availableTokens has [%d]", e, availableTokens)
			fmt.Println()
		}
	}

	fmt.Println(fillInterval)
}

func main() {
	bucket()
}
