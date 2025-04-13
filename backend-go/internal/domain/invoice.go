package domain

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusPending  Status = "pending"
	StatusApproved Status = "approved"
	StatusRejected Status = "rejected"
)

type InvoiceOptions struct {
	AccountID   string
	Amount      float64
	Description string
	PaymentType string
	Card        CreditCard
}

type Invoice struct {
	ID             string
	AccountID      string
	Amount         float64
	Status         Status
	Description    string
	PaymentType    string
	CardLastDigits string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type CreditCard struct {
	Number         string
	CVV            string
	ExpiryMonth    int
	ExpiryYear     int
	CardHolderName string
}

func NewInvoice(opts InvoiceOptions) (*Invoice, error) {
	if opts.Amount <= 0 {
		return nil, ErrInvalidAmount
	}

	if opts.AccountID == "" {
		return nil, ErrInvalidAccount
	}

	lastDigits := opts.Card.Number[len(opts.Card.Number)-4:]

	return &Invoice{
		ID:             uuid.New().String(),
		AccountID:      opts.AccountID,
		Amount:         opts.Amount,
		Status:         StatusPending,
		Description:    opts.Description,
		PaymentType:    opts.PaymentType,
		CardLastDigits: lastDigits,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}, nil
}

func (i *Invoice) Process() error {
	if i.Amount > 10000 {
		return nil
	}

	randomSource := rand.New(rand.NewSource(time.Now().Unix()))

	var newStatus Status
	if randomSource.Float64() <= 0.7 {
		newStatus = StatusApproved
	} else {
		newStatus = StatusRejected
	}

	i.Status = newStatus

	return nil
}

func (i *Invoice) UpdateStatus(newStatus Status) error {
	if i.Status != StatusPending {
		return ErrInvalidStatus
	}

	i.Status = newStatus
	i.UpdatedAt = time.Now()

	return nil
}
