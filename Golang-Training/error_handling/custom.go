package main

import (
	"errors"
	"fmt"
)

type CustomError struct {
	Data string
	Err error
}

func (c *CustomError) Error() string {
	return "custom error :" +c.Err.Error()
}

func (c *CustomError) Unwrap() error {
	return c.Err
}

var err error

func Custom() {

	err =  &CustomError{Data: "error data", Err: errors.New("internal error")}

	var customError *CustomError
	if errors.As(err, &customError) {
		fmt.Println("using errors.As, it is a custom error, data:", customError.Data)
	}

	if customError, ok := err.(*CustomError); ok {
		fmt.Println("using type assertion, it is a custom error, data:", customError.Data)
	}


	if errors.Is(returnCustomError(), err) {
		customError := err.(*CustomError)
		fmt.Println("using errors.Is, it is a custom error, data:", customError.Data)
	}
}

func returnCustomError() error {
	return err
}


