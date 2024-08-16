package user

import "errors"

var (
	ErrEmailAlreadyExists   = errors.New("email already exists")
	ErrAddressAlreadyExists = errors.New("address of this type already exists")
	ErrNotFound             = errors.New("user not found")
	ErrAddressNotFound      = errors.New("address not found")
)
