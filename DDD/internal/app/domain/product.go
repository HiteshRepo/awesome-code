package domain

const weightFactor = 0.01

type Product struct {
	name          string
	price         *Price
	weightInGrams int
}

func (p *Product) GetPrice() *Price {
	return p.price
}

func (p *Product) GetProductName() string {
	return p.name
}

func (p *Product) GetWeightInGrams() int {
	return p.weightInGrams
}

func (p *Product) GetCost() float64 {
	return p.price.GetValue() + (float64(p.weightInGrams) * weightFactor)
}

func NewProduct(name string, price *Price, weightInGrams int) Product {
	return Product{
		name:          name,
		price:         price,
		weightInGrams: weightInGrams,
	}
}
