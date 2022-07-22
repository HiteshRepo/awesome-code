package main

import "fmt"

func Allocations() {
	main1()
	main2()
	main3()
}

func main1() {
	n := 4
	n2 := square(n)
	fmt.Println(n2)
}

func square(x int) int {
	return x*x
}

func main2() {
	n := 4
	inc(&n)
	fmt.Println(n)
}

func inc(x *int) {
	*x++
}

func main3() {
	n := answer()
	fmt.Println(n)
}

func answer() *int {
	x := 42
	return &x
}
