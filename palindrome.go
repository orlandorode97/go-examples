package main

import (
	"fmt"
	"strings"
)

func main() {
	str := "doc note i dissent a fast never prevents a fatness i diet on cod"
	str = strings.ReplaceAll(str, " ", "")
	str = strings.ToLower(str)
	fmt.Println(isPalindrome(str))

}

func isPalindrome(str string) bool {
	middle := len(str) / 2

	for i, j := 0, len(str)-1; i < middle && j > middle; i, j = i+1, j-1 {
		if str[i] != str[j] {
			return false
		}
	}
	return true
}
