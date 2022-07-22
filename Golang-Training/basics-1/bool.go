package main

import "fmt"

func Bool() {
	str1 := "GolangTraining"
	str2 := "golangTraining"
	str3:= "GolangTraining"

	var result1 bool
	result1 = str1 == str2

	result2:= str1 == str3

	fmt.Println( result1)
	fmt.Println( result2)

	fmt.Printf("The type of result1 is %T and "+
		"the type of result2 is %T",
		result1, result2)
}
