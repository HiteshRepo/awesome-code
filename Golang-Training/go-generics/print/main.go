package main

import "fmt"

type ID int
type NewID int64

func main() {
	PrintThings(1, 2, 3, 4)
	// PrintThings(1.2, 2, 3, 4) --> will not work
	PrintThings(1.2, 2.5, 3, 4)
	PrintThings(1.2, 2.5, "a", 4)
	PrintThings(1, 2, "b", ID(4)) // working as underlying of ID is int itself
	// PrintThings(1, 2, "b", NewID(4)) -->  will not work as NEWID is not int rather int64, ~int denotes int or underlying should be int
}

func PrintThings[A, B any, C ~int](a1, a2 A, b B, c C) {
	fmt.Printf("%v, %v, %v, %v\n", a1, a2, b, c)
}
