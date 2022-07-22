package main

import "fmt"

func sum(nums ...int) {
	fmt.Print(nums, " ")
	total := 0
	for _, num := range nums {
		total += num
	}
	fmt.Println(total)
}

func displayUserInfo(info ...interface{}) {
	name := info[0].(string)
	age := info[1].(int)

	fmt.Println("Name of the user is", name)
	fmt.Println("Age of the user is", age)
}

func VariadicFunctions() {
	// sum of 2 numbers
	// other languages: sum(num1 int, num2 int)
	num1 := 1
	num2 := 2
	sum(num1, num2) // similar to params concept in other languages

	num3 := 3
	sum(num1, num2, num3)

	name := "Ramesh"
	age := 28
	displayUserInfo(name, age)
}


// sum(num1 int, num2 int, num3 int)
