package main

type s struct {
	a *string
	b *string
}

func main() {
	v := "Hello"

	str := s{a: &v, b: nil}
}

func toStringValue(st *string) string {
	if st == nil {
		return ""
	}
	return *st
}
