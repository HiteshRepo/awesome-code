package completed

type Driver struct{
	fizzBuzzCalculator FizzBuzzCalculator
}
type FizzBuzzCalculator interface {
	Calculate(number int) string
}
func (d *Driver)Drive(n int) string {
	result := ""
	for i:= 1; i < n ; i++ {
		result += d.fizzBuzzCalculator.Calculate(i) + ", "
	}
	result += d.fizzBuzzCalculator.Calculate(n)
	return result
}