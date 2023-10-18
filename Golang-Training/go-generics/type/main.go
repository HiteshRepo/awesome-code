package main

import "fmt"

func main() {
	intSet := NewSet[int](1, 2, 3)
	fmt.Printf("Has 2? %v\n", intSet.Has(2))
	fmt.Printf("Has 24? %v\n", intSet.Has(24))

	strSet := NewSet[string]("a", "b", "c")
	fmt.Printf("Has a? %v\n", strSet.Has("a"))
	fmt.Printf("Has d? %v\n", strSet.Has("d"))
}

type Set[T comparable] map[T]struct{}

func (s Set[T]) Has(value T) bool {
	_, ok := s[value]
	return ok
}

func NewSet[T comparable](values ...T) Set[T] {
	set := make(Set[T])
	for _, v := range values {
		set[v] = struct{}{}
	}

	return set
}
