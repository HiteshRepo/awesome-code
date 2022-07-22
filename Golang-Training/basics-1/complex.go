package main

import "fmt"

func Complex() {
	var a complex128 = complex(6, 2)
	var b complex64 = complex(9, 2)
	fmt.Println(a)
	fmt.Println(b)

	fmt.Printf("The type of a is %T and "+
		"the type of b is %T", a, b)
}
