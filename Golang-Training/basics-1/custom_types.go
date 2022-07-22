package main

import "fmt"

type actualPrice float64
type discountedPrice float64 // custom type

type aliasPrice = float64 // aliasing

type colors int

const (
	red colors = 1
	blue = 2
	white = 3
	green = 4
)

const (
	pi = 3.14
	si = 4.86
)

func CustomTypes() {
	ap := 29.32
	dp := 29.32 * 0.9

	displayPriceByANoob(ap, dp)

	displayPrice(dp, ap)

	ap2 := actualPrice(29.32)
	dp2 := discountedPrice(29.32 * 0.9)

	displayPriceBetter(ap2, dp2)

	displayColors(red)
	displayColors(blue)
}

func displayColors(inputColor colors) {
	fmt.Println(inputColor)
}

func displayPriceByANoob(price1, price2 float64) {
	fmt.Println("actual price of the ice cream is", price2)
	fmt.Println("discounted price of the ice cream is", price1)
}

func displayPrice(actualPrice, discountedPrice float64) {
	fmt.Println("actual price of the ice cream is", actualPrice)
	fmt.Println("discounted price of the ice cream is", discountedPrice)
}

func displayPriceBetter(actualPrice actualPrice, discountedPrice discountedPrice) {
	fmt.Println("actual price of the ice cream is", actualPrice)
	fmt.Println("discounted price of the ice cream is", discountedPrice)
}

//type dbDetails struct {
//	name
//	user
//	pass
//	port
//}
//
//type pgDetails dbDetails
//type nosqlDetails dbDetails
//
//func DBConn(dbDetails1 nosqlDetails, dpDetails2 pgDetails) {
//
//}
