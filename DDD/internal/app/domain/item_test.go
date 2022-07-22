package domain_test

import (
	"github.com/hiteshpattanayak-tw/DDD/internal/app/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewItem_ShouldReturnNewItem(t *testing.T) {
	product := domain.NewProduct("IPad", domain.NewPrice("INR", 60000), 1500)
	expected := domain.Item{Product: product, Qty: 1}
	actual := domain.NewItem(product, 1)

	assert.Equal(t, expected, actual)
}
