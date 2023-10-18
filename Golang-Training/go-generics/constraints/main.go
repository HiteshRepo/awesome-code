package main

import (
	"fmt"
	"math"
)

func main() {
	equalInt := func(a, b int) bool { return a == b }
	fmt.Printf("Has 1? %v\n", Has[int]([]int{1, 2, 3, 4}, 1, equalInt))

	equalStr := func(a, b string) bool { return a == b }
	fmt.Printf("Has a? %v\n", Has[string]([]string{"a", "b", "c"}, "a", equalStr))

	PrintBalance[INR](INR(250))
}

func Has[T any](list []T, value T, equalFn func(a, b T) bool) bool {
	for _, v := range list {
		// if v == value { // will not work because T is of any constraint type which is not comparable
		if equalFn(v, value) {
			return true
		}
	}

	return false
}

type Currency interface {
	~int | ~int64
	IS04127() string
	Decimal() int
}

func PrintBalance[T Currency](c T) {
	balance := float64(c) / math.Pow10(c.Decimal())
	fmt.Printf("%.*f %s\n", c.Decimal(), balance, c.IS04127())
}

type INR int64

func (c INR) IS04127() string {
	return "INR"
}

func (c INR) Decimal() int {
	return 2
}
