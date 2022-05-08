package string_manipulator

import (
	"fmt"
	"strings"
)

type Manipulator struct {}

func NewManipulator() Manipulator {
	return Manipulator{}
}

func (m Manipulator) StripWhitespace(datum string) string {
	return strings.TrimSpace(datum)
}

func (m Manipulator) ToUppercase(datum string) string {
	return strings.ToUpper(datum)
}

func (m Manipulator) Reverse(datum string) string {
	characters := []rune(datum)

	i := 0
	j := len(datum) - 1

	for i < j {
		m.swap(i, j, characters)
		i += 1
		j -= 1
	}

	return string(characters)
}

func (m Manipulator) swap(i, j int, arr []rune) {
	arr[i], arr[j] = arr[j], arr[i]
}

func (m Manipulator) Display(datum string) {
	fmt.Println(datum)
}
