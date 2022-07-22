package main

import "fmt"

func Float() {
	var x float32
	x = 1.23
	fmt.Println(x)

	x2 := float64(2.34)
	fmt.Println(x2)
}

func main() {
	Float()
}
