package models

import "errors"

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credential")
	ErrDuplidatedEmail    = errors.New("models: email already exists")
)
