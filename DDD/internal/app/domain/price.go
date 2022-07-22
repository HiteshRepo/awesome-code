package domain

type Price struct {
	currency string
	value    float64
}

func NewPrice(currency string, value float64) *Price {
	return &Price{
		currency: currency,
		value:    value,
	}
}

func (p *Price) GetCurrency() string {
	return p.currency
}

func (p *Price) GetValue() float64 {
	return p.value
}
