package domain

import "errors"

var (
	// Account errors
	ErrAccountNotFound    = errors.New("account not found")
	ErrDuplicatedApiKey   = errors.New("duplicated api key")
	ErrUnauthorizedAccess = errors.New("unauthorized access")

	// Invoice errors
	ErrInvalidAmount             = errors.New("invalid amount")
	ErrInvalidCard               = errors.New("invalid card")
	ErrInvalidAccount            = errors.New("invalid account")
	ErrInvalidStatus             = errors.New("invalid status")
	ErrInvoiceNotFound           = errors.New("invoice not found")
	ErrInvoicesByAccountNotFound = errors.New("invoices by account not found")
)
