package mocks

type FizzBuzzDummy struct {}

func (f *FizzBuzzDummy)Calculate(n int) string {
	if n == 1 {
		return "1"
	}
	if n == 2 {
		return "Fizz"
	}
	if n == 3 {
		return "Buzz"
	}
	if n == 4 {
		return "FizzBuzz"
	}
	return ""
}