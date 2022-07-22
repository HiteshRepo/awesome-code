package main

import "fmt"

func Pointers() {
	{
		a := 3
		b := a // 'b' is a copy of 'a' (all assignments are copy-operations)

		b++

		if a != b {
			fmt.Println("A and B are independent and is not equal")
		}
	}

	{
		a := 3
		b := &a // 'b' is the address of 'a'

		*b = *b + 2          // de-referencing 'b' means acting like a mutable copy of 'a'

		if a == *b {
			fmt.Println("A and B are same and is equal")
		}
	}

	{
		increment := func(i int) {
			i++
		}

		a := 3
		increment(a)

		if a == 3 {
			fmt.Println("A is still 3, as variables are always passed by value, and so a copy is made")
		}
	}

	{
		realIncrement := func(i *int) {
			*i++
		}

		b := 3
		realIncrement(&b)

		if b == 4 {
			fmt.Println("B is 4, as variables reference is passed")
		}
	}
}