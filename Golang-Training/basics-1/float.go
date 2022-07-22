package main

import "fmt"

func Float() {
	var x float32
	x = 1.23
	fmt.Println(x)

	a := 20.45
	b := 34.89

	c := b-a

	fmt.Printf("Result is: %f", c)
	fmt.Printf("\nThe type of c is : %T", c)
}
