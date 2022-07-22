package main

import "fmt"

func Enumeration() {
	states := map[string]string{
		"Maharashtra": "Mumbai",
		"Tamil Nadu": "Chennai",
		"West Bengal": "Kolkata",
		"Madhya Pradesh": "Bhopal",
	}

	for key, value := range states {
		fmt.Printf("state: %s, ", key)
		fmt.Printf("capital: %s\n", value)
	}
}
