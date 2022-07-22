package main

import "fmt"

type person struct {
	name string
	age  int
}

func(p person) structFunction() {
	fmt.Println("Hello from person struct function ", p.name, " and ", p.age)
}

func Structs() {
	p := person {
		name: "foo",
		age: 18,
	}

	p.structFunction()
}
