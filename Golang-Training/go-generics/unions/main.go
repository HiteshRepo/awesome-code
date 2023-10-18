package main

import "fmt"

func main() {
	fmt.Printf("%v compares %v? %v\n", 1, 1, Compare[int](1, 1))
	fmt.Printf("%v compares %v? %v\n", 1, 2, Compare[int](1, 2))

	fmt.Printf("%v compares %v? %v\n", true, false, Compare[bool](true, false))
}

func Compare[T int | bool](a, b T) int {
	if a == b {
		return 0
	}

	return 1
}
