// This file is for syntax testing and syntax cheat sheet only
package main

import (
	"fmt"
)

func main() {
	var i float32
	j := "mueed-dev"
	i = 47
	fmt.Printf("Hello, playground\n")
	fmt.Printf("%v, %T\n", i, i)
	fmt.Printf("%v\n", j)

	var numbers = make([]int, 3, 5)
	printSlice(numbers)
}
func printSlice(x []int) {
	fmt.Printf("len=%d cap=%d slice=%v\n", len(x), cap(x), x)
}
