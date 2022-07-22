package domain_test

import (
	"fmt"
	"github.com/hiteshpattanayak-tw/DDD/internal/app/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddCart_ShouldAddNewProductToCart(t *testing.T) {
	product := domain.NewProduct("IPad", domain.NewPrice("INR", 60000), 1500)
	item := domain.NewItem(product, 2)

	expected := &domain.Cart{Items: []domain.Item{item}}

	cart := domain.NewCart()
	cart.Add(item)

	assert.ElementsMatch(t, expected.Items, cart.Items)
}

func TestRemoveCart_ShouldRemoveProductFromCartAndAddItToRemovedItems(t *testing.T) {
	product1 := domain.NewProduct("IPad", domain.NewPrice("INR", 60000), 1500)
	product2 := domain.NewProduct("Pen", domain.NewPrice("INR", 10), 10)
	item1 := domain.NewItem(product1, 2)
	item2 := domain.NewItem(product2, 1)

	domainEvents := []string{
		fmt.Sprintf("%s was added", item1.GetProductName()),
		fmt.Sprintf("%s was added", item2.GetProductName()),
		fmt.Sprintf("%s was removed", item2.GetProductName()),
	}

	expected := &domain.Cart{Items: []domain.Item{item1}, DomainEvents: domainEvents}

	cart := domain.NewCart()
	cart.Add(item1)
	cart.Add(item2)

	cart.Remove(item2)

	assert.ElementsMatch(t, expected.Items, cart.Items)
	assert.ElementsMatch(t, expected.DomainEvents, cart.DomainEvents)
}

func TestIsEqualTo_ShouldReturnFalseForDifferentCarts(t *testing.T) {
	product := domain.NewProduct("IPad", domain.NewPrice("INR", 60000), 1500)
	item := domain.NewItem(product, 2)

	cart1 := domain.NewCart()
	cart2 := domain.NewCart()

	cart1.Add(item)
	cart2.Add(item)

	assert.False(t, cart1.IsEqualTo(cart2))
}

func TestIsEqualTo_ShouldReturnTrueForSameCarts(t *testing.T) {
	product := domain.NewProduct("IPad", domain.NewPrice("INR", 60000), 1500)
	item := domain.NewItem(product, 2)

	cart1 := domain.NewCart()
	cart1.Add(item)

	assert.True(t, cart1.IsEqualTo(cart1))
}
