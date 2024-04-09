package main

import (
	"fmt"
)

func main() {
	message := make(chan int, 3)

	message <- 5
	message <- 2
	message <- -100

	x := <-message
	fmt.Println(x)
	y := <-message
	fmt.Println(y)
	z := <-message
	fmt.Println(z)
}
