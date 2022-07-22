package main

import "fmt"

var acc int

func Defer() {
	increment := func(value int) {
		acc += value
		//fmt.Println("Increment is called")
	}

	decrement := func(value int) {
		acc -= value
		//fmt.Println("Decrement is called")
	}

	panicRecover := func() {
		if r := recover(); r != nil {
			fmt.Println("Panicked with the error ", r)
			increment(1)
		}
	}

	func() {
		acc = 0
		defer increment(1)
	}()
	if acc == 1 {
		fmt.Println("acc is 1 as defer function will be executed after main function body")
	}

	func() {
		acc = 0
		defer increment(5)
		defer decrement(3)
		if acc == 0 {
			fmt.Println("acc is still 0 as defer function will be executed only after main function body")
		}
	}()
	// Guess which increment gets called first :)

	func() {
		defer panicRecover()
		acc = 0
		panic("Expected error")
	}()
	if acc == 1 {
		fmt.Println("The program still continues as panic is recovered")
	}

}
