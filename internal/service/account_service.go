package service

import (
	"github.com/lauz1n/go-gateway/internal/domain"
	"github.com/lauz1n/go-gateway/internal/dto"
)

type AccountService struct {
	repository domain.AccountRepository
}

func NewAccountService(repository domain.AccountRepository) *AccountService {
	return &AccountService{repository: repository}
}

func (s *AccountService) CreateAccount(input *dto.CreateAccountInput) (*dto.AccountOutput, error) {
	account := dto.ToAccount(input)

	existingAccount, err := s.repository.FindByAPIKey(account.APIKey)
	if err != nil && err != domain.ErrAccountNotFound {
		return nil, err
	}
	if existingAccount != nil {
		return nil, domain.ErrDuplicatedApiKey
	}

	err = s.repository.Save(account)
	if err != nil {
		return nil, err
	}

	return dto.FromAccount(account), nil
}

func (s *AccountService) UpdateBalance(apiKey string, amount float64) (*dto.AccountOutput, error) {
	account, err := s.repository.FindByAPIKey(apiKey)
	if err != nil {
		return nil, err
	}

	account.AddBalance(amount)
	err = s.repository.UpdateBalance(account)

	if err != nil {
		return nil, err
	}

	return dto.FromAccount(account), nil
}

func (s *AccountService) FindByApiKey(apiKey string) (*dto.AccountOutput, error) {
	account, err := s.repository.FindByAPIKey(apiKey)
	if err != nil {
		return nil, err
	}

	return dto.FromAccount(account), nil
}

func (s *AccountService) FindById(id string) (*dto.AccountOutput, error) {
	account, err := s.repository.FindById(id)
	if err != nil {
		return nil, err
	}

	return dto.FromAccount(account), nil
}
