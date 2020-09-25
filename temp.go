// This file is for syntax testing and syntax cheat sheet only
package main

import (
	"fmt"
)

func main() {
	var i float32
	j := "moh-dev"
	i = 47
	fmt.Printf("Hello, playground\n")
	fmt.Printf("%v, %T\n", i, i)
	fmt.Printf("%v\n", j)

	var numbers = make([]int, 3, 5)
	printSlice(numbers)

	var fruitArr [2]string
	arr2 := [3]string{"App", "aaa", "Graaape"}
	fruitArr[0] = "Apple"

	fmt.Println(arr2[1:2])
	for i := 1; i <= 10; i++ {
		//	fmt.Print(i)
	}
	//Mapss
	emails := map[string]string{"Bob": "mai;@hhh", "Opo": "Opo@gmail.com"}
	fmt.Println(emails["Bob"])

	////Range
	ids := []int{22, 33, 4, 5, 2, 33, 55}

	for _, id := range ids {
		fmt.Printf("ID: %d", id)
	}
	//Range with map[]type
	for k, v := range emails {
		fmt.Printf("%s: %s\n", k, v)
	}
}
func printSlice(x []int) {
	//	fmt.Printf("len=%d cap=%d slice=%v\n", len(x), cap(x), x)
}
