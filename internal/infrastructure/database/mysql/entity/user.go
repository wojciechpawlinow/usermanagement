package entity

import (
	"gopkg.in/guregu/null.v4"
)

type DbUser struct {
	ID          null.Int    `db:"id" json:"id"`
	UUID        null.String `db:"uuid" json:"uuid"`
	Email       null.String `db:"email" json:"email"`
	FirstName   null.String `db:"first_name" json:"first_name"`
	LastName    null.String `db:"last_name" json:"last_name"`
	PhoneNumber null.String `db:"phone_number" json:"phone_number"`
}

type DbAddress struct {
	Type       null.Int    `db:"type" json:"type"`
	Street     null.String `db:"street" json:"street"`
	City       null.String `db:"city" json:"city"`
	State      null.String `db:"state" json:"state"`
	PostalCode null.String `db:"postal_code" json:"postal_code"`
	Country    null.String `db:"country" json:"country"`
}
