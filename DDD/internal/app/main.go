package main

import (
	"github.com/hiteshpattanayak-tw/DDD/internal/app/domain"
)

func main() {
	ipad := domain.NewProduct("IPad Air 2020", domain.NewPrice("INR", 60000))
	pen := domain.NewProduct("Hero ink pen", domain.NewPrice("INR", 10))
	bat := domain.NewProduct("GM Cricket bat", domain.NewPrice("INR", 2000))

	item1 := domain.NewItem(ipad, 1)
	item2 := domain.NewItem(pen, 1)
	item3 := domain.NewItem(bat, 2)

	cart := domain.NewCart()
	cart.Add(item1)
	cart.Add(item2)
	cart.Add(item3)
}
