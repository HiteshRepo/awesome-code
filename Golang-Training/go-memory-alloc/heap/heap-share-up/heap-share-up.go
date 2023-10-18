package heap

import "fmt"

func StartAnswer() {
	n := answer()
	fmt.Println(*n / 2)
}

func answer() *int {
	x := 42
	return &x
}
