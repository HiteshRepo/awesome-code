package completed

import (
	"github.com/hiteshpattanayak-tw/golang-training/tdd/completed/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Driver(t *testing.T) {
	t.Run("Given n is valid positive number return fizz buzz string", func(t *testing.T) {
		driver := Driver{fizzBuzzCalculator: &mocks.FizzBuzzDummy{}}

		fizzBuzzString := driver.Drive(4)

		assert.Equal(t,"1, Fizz, Buzz, FizzBuzz",
			fizzBuzzString)
	})

	t.Run("Given n is valid bigger positive number return fizz buzz string", func(t *testing.T) {
		fizzBuzzMock := mocks.FizzBuzzCalculator{}
		driver := Driver{fizzBuzzCalculator: &fizzBuzzMock}


		fizzBuzzMock.On("Calculate", 1).Return("1")
		fizzBuzzMock.On("Calculate", 2).Return("Fizz")
		fizzBuzzMock.On("Calculate", 3).Return("Buzz")
		fizzBuzzMock.On("Calculate", 4).Return("FizzBuzz")
		fizzBuzzMock.On("Calculate", 5).Return("BlahBlah")
		//fizzBuzzMock.On("Calculate", 6).Return("BlahBlah")

		fizzBuzzString := driver.Drive(5)

		assert.Equal(t,"1, Fizz, Buzz, FizzBuzz, BlahBlah",
			fizzBuzzString)
		fizzBuzzMock.AssertExpectations(t)
	})
}
