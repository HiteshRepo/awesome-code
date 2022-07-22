package main

import "fmt"

func RemoveDups(nums []int) int {
	i := 0
	for j:=1; j<len(nums); j++ {
		if nums[i] != nums[j] {
			i += 1
			nums[i] = nums[j]
		}
	}
	return i+1
}

func main() {
	arr := []int{0,0,1,1,1,2,2,3,3,4}
	ans := RemoveDups(arr)
	fmt.Println(arr)
	fmt.Println(arr[:ans])
}
