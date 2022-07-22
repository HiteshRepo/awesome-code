package main

import "fmt"

func AnonymousFunctions() {

	func() {
		fmt.Println("Hello from AnonymousFunctions - Self Calling")
	}()

	storedFn := func() {
		fmt.Println("Hello from AnonymousFunctions - Stored")
	}
	storedFn()

	func(entityName string) {
		fmt.Println("Hello from AnonymousFunctions -", entityName)
	}("Parameterized")

	returnAnonymousFunction()("Returned Parameterized")
}

func returnAnonymousFunction() func(string) {
	return func(entityName string) {
		fmt.Println("Hello from AnonymousFunctions -", entityName)
	}
}
