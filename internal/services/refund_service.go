package services

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/constants"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
)

type RefundService interface {
	CronAutoApproveRefundByExpirationTime() error
}

type refundService struct {
	refundRepository r.RefundRepository
}

type RefundServiceConfig struct {
	RefundRepository r.RefundRepository
}

func NewRefundService(c *RefundServiceConfig) RefundService {
	return &refundService{
		refundRepository: c.RefundRepository,
	}
}

func (p *refundService) CronAutoApproveRefundByExpirationTime() error {
	newRefunds := &models.Refund{
		Status: constants.RefundStatusApprovedBySystem,
	}
	err := p.refundRepository.UpdateAllByStatusByRefundExpirationDatetimeLteNow(constants.RefundStatusOngoing, newRefunds)
	if err != nil {
		return err
	}

	return nil
}
