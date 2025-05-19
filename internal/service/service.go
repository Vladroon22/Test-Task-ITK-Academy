package service

import (
	"context"

	"github.com/Vladroon22/TestTask-ITK-Academy/internal/entity"
)

type Servicer interface {
	WalletOperation(c context.Context, wallet entity.WalletData) error
	GetBalance(c context.Context, uuid string) (entity.WalletData, error)
}

type Service struct {
	repo Servicer
}

func NewService(repo Servicer) Servicer {
	return &Service{repo: repo}
}

func (s *Service) GetBalance(c context.Context, uuid string) (entity.WalletData, error) {
	return s.repo.GetBalance(c, uuid)
}

func (s *Service) WalletOperation(c context.Context, wallet entity.WalletData) error {
	return s.repo.WalletOperation(c, wallet)
}
