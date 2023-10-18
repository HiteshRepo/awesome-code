package main

import "fmt"

func main() {

	fmt.Printf("Has a? %v\n", Has[string]([]string{"a", "b"}, "a"))
	fmt.Printf("Has c? %v\n", Has[string]([]string{"a", "b"}, "c"))
	fmt.Printf("Has 1? %v\n", Has[int]([]int{1, 2}, 1))
	fmt.Printf("Has 3? %v\n", Has[int]([]int{1, 2}, 3))
}

func Has[T comparable](list []T, value T) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}

	return false
}
