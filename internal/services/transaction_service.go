package services

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/constants"
	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	externalDtos "git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	"github.com/google/uuid"
)

type TransactionService interface {
	PendingPaymentGetAllByUser(id uuid.UUID) ([]*externalDtos.PendingPaymentGetAllByUserResponse, error)
}

type transactionService struct {
	transactionRepository                  r.TransactionRepository
	sealabsPayTransactionPaymentRepository r.SealabsPayTransactionPaymentRepository
}

type TransactionServiceConfig struct {
	TransactionRepository                  r.TransactionRepository
	SealabsPayTransactionPaymentRepository r.SealabsPayTransactionPaymentRepository
}

func NewTransactionService(c *TransactionServiceConfig) TransactionService {
	return &transactionService{
		transactionRepository:                  c.TransactionRepository,
		sealabsPayTransactionPaymentRepository: c.SealabsPayTransactionPaymentRepository,
	}
}

func (s *transactionService) PendingPaymentGetAllByUser(id uuid.UUID) ([]*externalDtos.PendingPaymentGetAllByUserResponse, error) {
	transaction, err := s.transactionRepository.PendingPaymentGetAllByUserID(id)

	if err != nil {
		return nil, err
	}

	var pendingTransaction []*dtos.PendingPaymentTransaction

	for _, tx := range transaction {
		if tx.PaymentMethodID == constants.SealabsPayID {
			sealabsPay, err := s.sealabsPayTransactionPaymentRepository.GetByTransactionID(tx.TransactionID)

			if err != nil {
				return nil, err
			}

			tx.CardNumber = sealabsPay.SealabsPayCardNumber
		}

		pendingTransaction = append(pendingTransaction, tx)

	}

	return externalDtos.NewPendingPaymentGetAllByUserResponse(pendingTransaction), nil
}
