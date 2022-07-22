package domain

type Customer struct {
	address    *Address
	accountIds []string
}

func NewCustomer(address *Address) *Customer {
	return &Customer{
		address:    address,
		accountIds: make([]string, 0),
	}
}

func (c *Customer) AddAccount(accountId string) {
	c.accountIds = append(c.accountIds, accountId)
}

func (c *Customer) UpdateAddresses(address *Address) {
	c.address = address
}

func (c *Customer) GetAddress() *Address {
	return c.address
}
