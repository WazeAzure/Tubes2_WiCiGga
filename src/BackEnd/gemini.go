package main

import (
	"fmt"
)

func main() {
	a := []int{}

	for i := 0; i < 100000; i++ {
		a = append(a, i)
		if i%4 == 0 {
			a = a[1:]
		}
		fmt.Println(len(a), cap(a))
	}

	fmt.Println("hello world")
}
