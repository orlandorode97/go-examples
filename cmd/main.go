package main

import (
	"errors"
	"fmt"
)

// Define the Classification type here.

type Classification string

const (
	ClassificationDeficient = "ClassificationDeficient"
	ClassificationPerfect   = "ClassificationPerfect"
	ClassificationAbundant  = "ClassificationAbundant"
)

var ErrOnlyPositive = errors.New("only positive numbers")

func Classify(n int64) (Classification, error) {
	if n <= 0 {
		return "", ErrOnlyPositive
	}
	haveSeen := map[int64]bool{}
	haveSeen[1] = true
	for i := int64(2); i < n; i++ {
		divide := n / i
		if n%i == 0 {
			if _, ok := haveSeen[divide]; !ok {
				haveSeen[divide] = true
			}
		}

	}

	sum := int64(0)
	for number := range haveSeen {
		sum += number
	}
	fmt.Printf("Debug message: sum %v and n %v\n ", sum, n)
	switch {
	case sum < n:
		return ClassificationDeficient, nil
	case sum == n:
		return ClassificationPerfect, nil
	default:
		return ClassificationAbundant, nil
	}
}

func main() {
	// Example usage
	classification, err := Classify(28)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Classification of 28 is:", classification)
	}
}
