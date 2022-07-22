package main

import (
	"fmt"
	"github.com/hiteshpattanayak-tw/golang-training/packages/numbers"
	"github.com/hiteshpattanayak-tw/golang-training/packages/strings"
)

func main() {
	inpStr := "Hello"
	fmt.Printf("length of string %s is %d\n", inpStr, strings.Length(inpStr))

	num1 := 12
	num2 := 23
	fmt.Printf("greater number between %d and %d is %d\n", num1, num2, numbers.Greater(num1, num2))
}
