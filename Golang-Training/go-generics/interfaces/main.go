package main

import "fmt"

func main() {
	fmt.Printf("Has 1? %v\n", Has[ID]([]ID{1, 2, 3}, 1))
}

type Equalizer[T any] interface {
	Equal(other T) bool
}

func Has[T Equalizer[T]](list []T, value T) bool {
	for _, v := range list {
		if v.Equal(value) {
			return true
		}
	}
	return false
}

type ID int

func (id ID) Equal(otherId ID) bool {
	return id == otherId
}
