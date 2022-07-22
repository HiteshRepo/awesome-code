package domain_service_test

import (
	"github.com/hiteshpattanayak-tw/DDD/internal/app/domain"
	"github.com/hiteshpattanayak-tw/DDD/internal/app/domain/domain_service"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOrderService_Checkout(t *testing.T) {
	product := domain.NewProduct("IPad", domain.NewPrice("INR", 60000), 1500)
	item := domain.NewItem(product, 2)

	cart := domain.NewCart()
	cart.Add(item)

	products := []domain.Product{
		product, product,
	}
	expected := domain.NewOrder(products)

	os := domain_service.NewOrderService()
	actual := os.Checkout(cart)
	assert.Equal(t, expected, actual)
	assert.True(t, cart.IsCheckedOut())
	assert.Empty(t, cart.Items)
}

func TestOrderService_CalculateTotalCost(t *testing.T) {
	product1 := domain.NewProduct("IPad", domain.NewPrice("INR", 60000), 1500)
	product2 := domain.NewProduct("Pen", domain.NewPrice("INR", 100), 100)
	item1 := domain.NewItem(product1, 2)
	item2 := domain.NewItem(product2, 10)

	cart := domain.NewCart()
	cart.Add(item1)
	cart.Add(item2)

	expectedCost := float64(121040) // ((60000 * 2) + (1500 * 0.01 * 2)) + ((100 * 2) + (100 * 0.01 * 10))

	os := domain_service.NewOrderService()
	actualCost := os.CalculateTotalCost(cart)

	assert.Equal(t, expectedCost, actualCost)
}
