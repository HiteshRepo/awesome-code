package main

import (
	"fmt"
)

func Slices() {
	// Integer Slice
	intSlice1 := []int{1,2,3,4,5}

	var intSlice2 []int
	intSlice2 = []int{1,2,3,4,5}

	intSlice3 := make([]int, 5)
	intSlice3 = append(intSlice3, 1)
	intSlice3 = append(intSlice3, 2)
	intSlice3 = append(intSlice3, 3)
	intSlice3 = append(intSlice3, 4)
	intSlice3 = append(intSlice3, 5)

	fmt.Println("Length of intSlice1", len(intSlice1))
	fmt.Println("Length of intSlice1", len(intSlice2))
	fmt.Println("Capacity of intSlice1", cap(intSlice1))
	fmt.Println("Capacity of intSlice1", cap(intSlice2))

	fmt.Println("Are intSlice1 and intSlice2 equal?", areArraysEqual(intSlice1, intSlice2))
	fmt.Println("Are intSlice1 and intSlice2 equal?", areArraysEqual(intSlice1, intSlice3))


	arr3 := []int{1,2,3,4,5}
	fmt.Println(cap(arr3))
	arr3 = append(arr3, 6)
	fmt.Println(cap(arr3))

	arr3 = append(arr3, 7)
	fmt.Println(cap(arr3))
	arr3 = append(arr3, 8,9,10)
	fmt.Println(cap(arr3))
}


func changeArr1(arr []int) {
	for i:=0; i<len(arr); i++ {
		arr[i] = arr[i]+1
	}
	arr = append(arr, 1000)
}

func changeArr2(arr *[]int) {
	for i:=0; i<len(*arr); i++ {
		(*arr)[i] = (*arr)[i]+1
	}
	*arr = append(*arr, 1000)
}

func changeArr3(arr [4]int) {
	for i:=0; i<len(arr); i++ {
		arr[i] = arr[i]+1
	}
}

func changeArr4(arr *[4]int) {
	for i:=0; i<len(*arr); i++ {
		(*arr)[i] = (*arr)[i]+1
	}
}


func areArraysEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}