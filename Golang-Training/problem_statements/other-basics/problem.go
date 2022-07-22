package main

import "fmt"

type Solution interface {
	GetDivideFunction() func(...int) int
}

func Problem(solution Solution)  {
	divide := solution.GetDivideFunction()

	val := divide(10, 2, 5)
	if val != 1{
		panic("wrong answer the answer must be 1")
	}

	val = divide(10, 0)
	if val != 0 {
		panic("wrong answer the answer must be 0")
	}

	fmt.Println("All answers are correct")
}