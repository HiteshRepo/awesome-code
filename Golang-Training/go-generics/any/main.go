package main

import "fmt"

func main() {
	intList := NewEmptyList[int]()
	intList = append(intList, 1)
	intList = append(intList, 2)
	fmt.Printf("Integer list: %v\n", intList)

	strList := NewEmptyList[string]()
	strList = append(strList, "a")
	strList = append(strList, "b")
	fmt.Printf("String list: %v\n", strList)
}

func NewEmptyList[T any]() []T {
	return make([]T, 0)
}
