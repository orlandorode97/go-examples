package main

import "fmt"

func main() {

	arr := []int{1, 2, 3}
	newArr := make([]int, 0, len(arr))
	for _, v := range arr {
		newArr = append(newArr, v)
	}
	fmt.Println(newArr, len(newArr), cap(newArr))
}
