package domain_service

import "github.com/hiteshpattanayak-tw/DDD/internal/app/domain"

const discount = 0.1

var CompetitorsProductList map[string]domain.Price

type CompetitorsBasedPricer struct{}

func NewCompetitorsBasedPricer() CompetitorsBasedPricerService {
	return &CompetitorsBasedPricer{}
}

func init() {
	CompetitorsProductList = make(map[string]domain.Price)
}

func (c CompetitorsBasedPricer) AddNewProductToCompetitorsList(product string, price domain.Price) {
	CompetitorsProductList[product] = price
}

func (c CompetitorsBasedPricer) GetDiscountedPrice(productName string) domain.Price {
	for prodName, price := range CompetitorsProductList {
		if prodName != productName {
			continue
		}
		val := price.GetValue()
		price := domain.NewPrice("INR", val*(1-discount))
		return *price
	}

	return domain.Price{}
}
