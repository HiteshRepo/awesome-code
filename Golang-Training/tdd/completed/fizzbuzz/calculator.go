package fizzbuzz

import "strconv"

func Calculate(number int) string {
	result := ""
	if number % 3 == 0 {
		result += "Fizz"
	}
	if number % 5 == 0 {
		result += "Buzz"
	}
	if result != "" {
		return result
	}
	return strconv.Itoa(number)
}