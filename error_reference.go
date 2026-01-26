package main

import (
	"errors"
	"fmt"
)

func main() {

	err := errors.New("this is my error")
	printError(&err)
	fmt.Println(err)
}

func printError(err *error) {
	fmt.Println(*err)
	a := errors.New("this pedio")
	&err = a
}
