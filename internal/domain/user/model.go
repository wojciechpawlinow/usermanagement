package user

import (
	"github.com/wojciechpawlinow/usermanagement/internal/domain"
)

type User struct {
	ID          domain.ID  `json:"id"`
	Email       string     `json:"email"`
	Password    string     `json:"-"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	PhoneNumber string     `json:"phone_number"`
	Addresses   []*Address `json:"addresses"`
}

type AddressType int

const (
	WorkAddress AddressType = iota
	HomeAddress
	BillingAddress
)

type Address struct {
	Type       AddressType `json:"type"`
	Street     string      `json:"street"`
	City       string      `json:"city"`
	State      string      `json:"state"`
	PostalCode string      `json:"postal_code"`
	Country    string      `json:"country"`
}
