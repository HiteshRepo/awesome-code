package strings

func Length(inputString string) int {
	count := 0
	for _,c := range inputString {
		_ = c
		count += 1
	}
	return count
}
