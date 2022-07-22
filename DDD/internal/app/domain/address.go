package domain

type Address struct {
	city string
}

func NewAddress(cityName string) *Address {
	return &Address{city: cityName}
}
