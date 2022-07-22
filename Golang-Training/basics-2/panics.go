package main

import (
	"fmt"
	"sync"
)

func Panics() {
	//PanicByProgram()
	//PanicByProgram2()
	//PanicByProgrammer()
	DeferAndPanic()
}

func PanicByProgram() {
	a := 0
	fmt.Println(1/a)
}

func PanicByProgram2() {
	numbers := []int{1,2,3,0}
	var wg sync.WaitGroup
	for _,n := range numbers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println(100/n)
		}()
	}

	wg.Wait()
}

func PanicByProgrammer() {
	empName := "Ramesh"
	age := 75

	employee(empName, age)
}

func employee(name string, age int){
	if age > 65{
		panic("Age cannot be greater than retirement age")
	}
	fmt.Printf("%s is within retirement age\n", name)
}

func DeferAndPanic() {
	A_lang := "GO Language"
	A_name := "test name"
	defer fmt.Println("Defer statement in the Calling function")
	// entry(&A_lang, nil)
	entry(&A_lang, &A_name)
}

func entry(lang *string, aname *string) {

	defer fmt.Println("Defer statement in the entry function")

	if lang == nil {
		panic("Error: Language cannot be nil")
	}

	if aname == nil {
		panic("Error: Author name cannot be nil")
	}

	fmt.Printf("Author Language: %s \n Author Name: %s\n", *lang, *aname)
}
