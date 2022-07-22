package domain_test

import (
	"github.com/hiteshpattanayak-tw/DDD/internal/app/domain"
	"github.com/hiteshpattanayak-tw/DDD/internal/app/domain/domain_service"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewProduct_ShouldReturnNewProduct(t *testing.T) {
	name := "pen"
	currency := "INR"
	value := 20.1
	price := domain.NewPrice(currency, value)

	expected := domain.NewProduct(name, price, 10)
	actual := domain.NewProduct(name, price, 10)

	assert.Equal(t, expected, actual)

}

func TestNewProduct_ShouldReturnDiscountedPriceIfProductPresent(t *testing.T) {
	name := "pen"
	currency := "INR"
	value := 20.1
	price := domain.NewPrice(currency, value)

	cbp := domain_service.NewCompetitorsBasedPricer()
	cbp.AddNewProductToCompetitorsList(name, *price)
	discountedPrice := cbp.GetDiscountedPrice(name)

	expected := domain.NewProduct(name, &discountedPrice, 10)
	actual := domain.NewProduct(name, &discountedPrice, 10)

	assert.Equal(t, expected, actual)
}
