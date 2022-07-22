package domain_service

import "github.com/hiteshpattanayak-tw/DDD/internal/app/domain"

type OrderService struct{}

func NewOrderService() CheckoutService {
	return &OrderService{}
}

func (o OrderService) Checkout(cart *domain.Cart) *domain.Order {
	order := domain.NewOrder(o.getFlattenedProducts(cart))
	cart.Checkout()
	return order
}

func (o OrderService) CalculateTotalCost(cart *domain.Cart) float64 {
	products := o.getFlattenedProducts(cart)
	cost := float64(0)
	for _, p := range products {
		cost += p.GetCost()
	}
	return cost
}

func (o OrderService) getFlattenedProducts(c *domain.Cart) []domain.Product {
	products := make([]domain.Product, 0)
	for _, item := range c.Items {
		i := 0
		for i < item.Qty {
			products = append(products, item.Product)
			i += 1
		}
	}

	return products
}
