package main

import "fmt"

func Maps() {
	ages := map[string]int{
		"bob": 10,
		"joe": 20,
		"dan": 30,
	}

	age := ages["bob"]
	fmt.Println("Maps: fetched age for bob ", age)

	age, ok := ages["steven"]
	if !ok {
		fmt.Println("Maps: stevens' age is not present ")
	}

	fmt.Println("Maps: setting bob's age to 99")
	ages["bob"] = 99

	fmt.Println("Maps: deleting bob's age")
	delete(ages, "bob")
}
