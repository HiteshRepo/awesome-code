package domain

type Order struct {
	products []Product
}

func NewOrder(products []Product) *Order {
	return &Order{products: products}
}

func (o *Order) GetProducts() []Product {
	return o.products
}
