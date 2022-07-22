package fizzbuzz

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Calculate_ForNumberCases(t *testing.T) {
	t.Run("Given 1 Should return 1", func(t *testing.T) {
		number := 1
		actualNumber := Calculate(number)

		assert.Equal(t, "1", actualNumber)
	})

	t.Run("Given 8 Should return 8", func(t *testing.T) {
		number := 8
		actualNumber := Calculate(number)

		assert.Equal(t, "8", actualNumber)
	})

}


func Test_Calculate_ForFizz(t *testing.T) {
	t.Run("Given 3 Should return Fizz", func(t *testing.T) {
		number := 3
		actualNumber := Calculate(number)

		assert.Equal(t, "Fizz", actualNumber)
	})

	t.Run("Given 6 Should return Fizz", func(t *testing.T) {
		number := 6
		actualNumber := Calculate(number)

		assert.Equal(t, "Fizz", actualNumber)
	})
}

func Test_Calculate_ForBuzz(t *testing.T) {
	t.Run("Given 5 Should return Buzz", func(t *testing.T) {
		number := 5
		actualNumber := Calculate(number)

		assert.Equal(t, "Buzz", actualNumber)
	})

	t.Run("Given 10 Should return Buzz", func(t *testing.T) {
		number := 10
		actualNumber := Calculate(number)

		assert.Equal(t, "Buzz", actualNumber)
	})
}

func Test_Calculate_ForFizzBuzz(t *testing.T) {
	t.Run("Given 15 Should return FizzBuzz", func(t *testing.T) {
		number := 15
		actualNumber := Calculate(number)

		assert.Equal(t, "FizzBuzz", actualNumber)
	})
}