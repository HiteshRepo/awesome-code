package domain_test

import (
	"github.com/hiteshpattanayak-tw/DDD/internal/app/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddCart_ShouldReturnNewPrice(t *testing.T) {
	currency := "INR"
	value := 20.1

	actual := domain.NewPrice(currency, value)

	assert.Equal(t, currency, actual.GetCurrency())
	assert.Equal(t, value, actual.GetValue())
}
