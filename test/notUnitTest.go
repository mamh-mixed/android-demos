package main

import (
	"fmt"
	"time"
)

func main() {

	for i := 0; i < 100; i++ {
		go func(int) {
			fmt.Println(i)
		}(i)
	}

	time.Sleep(3 * time.Second)
}
