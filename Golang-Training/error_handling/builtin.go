package main

import (
	"errors"
	"fmt"
)

// Unconventional

// Built in interface type


func BuiltIn() {

	// Creation of error string types

	error1 := errors.New("this is a built in error created by errors.new")

	fmt.Println(error1.Error())

	// wrapping

	error2 := fmt.Errorf("this is a built in error created by fmt and wraps %w", error1)

	// unwrapping

	fmt.Println(errors.Unwrap(error2))  // error1.Error() and errors.Unwrap(error2) is equal

	// function returning error
	err := returnError()

	fmt.Println(err.Error())

}

func returnError() error {
	if false {
		return nil // can also be nil
	}
	return errors.New("the function returned error")
}