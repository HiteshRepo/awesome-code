package main

import (
	"github.com/thoas/go-funk"
	"log"
)

var validdata = []int{1,2,3,4,5}

func main() {
	reqData := []int{1,2,3,4,5}
	log.Println(len(funk.Intersect(validdata, reqData).([]int)) == len(validdata))
}
