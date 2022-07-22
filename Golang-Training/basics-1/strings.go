package main

import (
	"fmt"
	"strings"
)

func String() {
	//str := "GolangTraining"

	var str2 string
	str2 = "New golang training"

	//fmt.Printf("Name of course is %s and it costs 0 INR\n", str2)
	//fmt.Println("Name of course is", str2, "and it costs 0 INR")
	//fmt.Println(fmt.Sprintf("Name of course is %s and it costs 0 INR\n", str2))
	//
	fmt.Printf("Length of the string is:%d", len(str2))
	//
	//fmt.Printf("\nString is: %s", str)
	//
	//fmt.Printf("\nType of str is: %T", str)
	//
	//fmt.Printf("\nString2 is: %s", str2)
	//fmt.Printf("\nLength of String2 is: %d", len(str2))
	//
	//
	fmt.Printf("\nWords of String2 is: %d\n", len(strings.Split(str2, " ")))
	//
	//str3 := `String3 contains another string in double quotes: "ABC"`
	//fmt.Println(str3)
	//
	//name := "Ramesh"
	//fmt.Println("First character:", name[0]) // rune
	//fmt.Printf("Second character: %c\n", name[1])
	//fmt.Printf("Ascii of second character: %d\n", int(name[1]))
}
