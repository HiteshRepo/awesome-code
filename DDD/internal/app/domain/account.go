package domain

import gonanoid "github.com/matoous/go-nanoid/v2"

type Account struct {
	id      string
	address *Address
}

func NewAccount(address *Address) *Account {
	id, _ := gonanoid.New()
	return &Account{id: id, address: address}
}

func (a *Account) UpdateAddress(address *Address) {
	a.address = address
}

func (a *Account) GetAddress() *Address {
	return a.address
}

func (a *Account) GetAccountId() string {
	return a.id
}
