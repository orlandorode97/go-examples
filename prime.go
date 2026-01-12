package main

import "fmt"

// Define the Classification type here.

type Classification int

const (
	ClassificationDeficient Classification = iota
	ClassificationPerfect
	ClassificationAbundant
)

func Classify(n int64) (Classification, error) {
	haveSeen := make(map[int64]bool)

	for i := int64(2); i < n; i++ {
		divide := n / i
		if _, ok := haveSeen[divide]; !ok {
			haveSeen[divide] = true
		}
	}

	sum := int64(0)
	for number := range haveSeen {
		sum += number
	}
	return ClassificationAbundant, nil
}

func main() {
	fmt.Println(Classify(24))
}
