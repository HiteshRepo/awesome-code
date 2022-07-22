package main

import "fmt"

func Enumeration() {
	SlicesEnumeration()
	ArrayEnumeration()
	StringEnumeration()
}

func SlicesEnumeration() {
	arr := []int{101, 102, 103, 104, 105}

	for i:=0; i<len(arr); i++ {
		fmt.Println(arr[i])
	}

	for i,n := range arr {
		fmt.Printf("index: %d, number: %d\n", i, n)
	}
}

func ArrayEnumeration() {
	arr := [5]int{101, 102, 103, 104, 105}

	for i:=0; i<len(arr); i++ {
		fmt.Println(arr[i])
	}

	for i,n := range arr {
		fmt.Printf("index: %d, number: %d\n", i, n)
	}
}

func StringEnumeration() {
	str := "String"

	for i:=0; i<len(str); i++ {
		fmt.Println(str[i])
	}

	for i,c := range str {
		fmt.Printf("index: %d, char: %c\n", i, c)
	}
}
