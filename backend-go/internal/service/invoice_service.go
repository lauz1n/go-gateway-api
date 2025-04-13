package service

import (
	"errors"
	"log"

	"github.com/lauz1n/go-gateway/internal/domain"
	"github.com/lauz1n/go-gateway/internal/dto"
)

type InvoiceService struct {
	invoiceRepository domain.InvoiceRepository
	accountService    *AccountService
}

func NewInvoiceService(repository domain.InvoiceRepository, accountService *AccountService) *InvoiceService {
	return &InvoiceService{invoiceRepository: repository, accountService: accountService}
}

func (s *InvoiceService) Create(input *dto.CreateInvoiceInput) (*dto.CreateInvoiceOutput, error) {
	accountOutput, err := s.accountService.FindByApiKey(input.APIKey)
	if err != nil {
		return nil, err
	}

	invoice, err := dto.ToInvoice(input, accountOutput.ID)
	if err != nil {
		return nil, err
	}

	if err := invoice.Process(); err != nil {
		return nil, err
	}

	if invoice.Status == domain.StatusApproved {
		s.accountService.UpdateBalance(input.APIKey, invoice.Amount)
		if err != nil {
			return nil, err
		}
	}

	if err := s.invoiceRepository.Save(invoice); err != nil {
		return nil, err
	}

	return dto.FromInvoice(invoice), nil
}

func (s *InvoiceService) GetById(id, apiKey string) (*dto.CreateInvoiceOutput, error) {
	invoice, err := s.invoiceRepository.FindById(id)
	if err != nil {
		return nil, err
	}

	accountOutput, err := s.accountService.FindByApiKey(apiKey)
	if err != nil {
		return nil, err
	}

	if invoice.AccountID != accountOutput.ID {
		log.Println("invoice account id and account output id are different", invoice.AccountID, accountOutput.ID)
		return nil, errors.New("invoice not found")
	}

	return dto.FromInvoice(invoice), nil
}

func (s *InvoiceService) ListByAccount(accountId string) ([]*dto.CreateInvoiceOutput, error) {
	invoices, err := s.invoiceRepository.FindByAccountID(accountId)
	if err != nil {
		return nil, err
	}

	output := make([]*dto.CreateInvoiceOutput, len(invoices))
	for i, invoice := range invoices {
		output[i] = dto.FromInvoice(invoice)
	}

	return output, nil
}

func (s *InvoiceService) ListByAccountApiKey(apiKey string) ([]*dto.CreateInvoiceOutput, error) {
	accountOutput, err := s.accountService.FindByApiKey(apiKey)
	if err != nil {

	}

	return s.ListByAccount(accountOutput.ID)
}
