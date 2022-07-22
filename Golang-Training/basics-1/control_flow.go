package main

import (
	"fmt"
	"math/rand"
	"time"
)

// if-else
// for
// switch-case

func ControlFlow() {
	//IfElse()
	//SwitchCase()
	//For()
	//BreakAndContinue()
}

func IfElse() {
	n := rand.Int()

	if n%2 == 0 {
		fmt.Printf("%d is an even number\n", n)
	} else {
		fmt.Printf("%d is an odd number\n", n)
	}

	if m := rand.Int(); m%2 == 0 {
		fmt.Printf("%d is an even number\n", m)
	} else {
		fmt.Printf("%d is an odd number\n", m)
	}

	// m is not accessible here, scope is within those if-else block only
	// n is accessible
}

func For() {
	// normal for loop syntax
	//for InitSimpleStatement; Condition; PostSimpleStatement {
	// do something
	//}

	for i := 0; i < 10; i++ {
		fmt.Println(i)
	}

	// behaves like a while loop
	var i = 0 // i := 0
	for i < 10 {
		fmt.Println(i)
		i++
	}

	for i := 0; ; i++ {
		if i >= 10 {
			break
		}
		fmt.Println(i)
	}
}

func SwitchCase() {
	// syntax

	//switch InitSimpleStatement; CompareOperand0 {
	//case CompareOperandList1:
	// do something
	//case CompareOperandList2:
	// do something
	//case CompareOperandListN:
	// do something
	//default:
	// do something
	//}

	// no break statement required at the end of each block
	// fallthrough is a keyword that can be used to let control pass through next case as well
	// arbitrage default position

	switch n := rand.Intn(3); n {
	case 0:
		fmt.Println("0")
	case 1:
		fmt.Println("1")
	default:
		fmt.Println("2")
	}

	switch h := time.Now().Hour(); {
	case h < 12:
		fmt.Println("Morning")
	case h > 12:
		fmt.Println("After 12")
	default:
		fmt.Println("This will never happen")
	}

	switch h := time.Now().Hour(); h < 12 {
	case true:
		fmt.Println("Before 12")
	case false:
		fmt.Println("After 12")
	default:
		fmt.Println("This will never happen")
	}

	switch h := time.Now().Hour(); h {
	case 5, 6, 7, 8, 9, 10, 11:
		fmt.Println("Morning")
	case 12, 13, 14, 15:
		fmt.Println("Afternoon")
	case 16, 17, 18, 19, 20:
		fmt.Println("Evening")
	default:
		fmt.Println("Night")
	}

	switch h := time.Now().Hour(); {
	case h < 12:
		fmt.Println("Morning")
		fallthrough
	case h < 11:
		fmt.Println("Still Morning")
		fallthrough
	case h < 8 && h > 5:
		fmt.Println("Early Morning")
	case h > 12:
		fmt.Println("Afternoon")
		fallthrough
	case h > 14 && h < 16:
		fmt.Println("Late Afternoon")
	case h > 16 && h < 20:
		fmt.Println("Evening")
	default:
		fmt.Println("Night")
	}
}

func BreakAndContinue() {
	// Prime Number above n

	n := 13
	found := false
	ans := -1
	for j:=n+1; ; j++{
		for i := 2; ; i++ {
			if i * i > j {
				ans = j
				found = true
				break
			}
			if j % i == 0 {
				break
			}
		}

		if found {
			break
		}
	}

	fmt.Println("ans = ", ans)



	ans = -1
	outer:
	for j:=n+1; ; j++{
		for i := 2; ; i++ {
			if i * i > j {
				ans = j
				break outer
			}
			if j % i == 0 {
				continue outer
			}
		}
	}

	fmt.Println("ans = ", ans)
}
