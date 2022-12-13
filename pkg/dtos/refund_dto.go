package dtos

import (
	"mime/multipart"
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	"github.com/google/uuid"
)

type OrderDetailRefund struct {
	OrderDetailID uuid.UUID `json:"order_detail_id" binding:"required"`
	Quantity      uint      `json:"quantity" binding:"required"`
}

type CreateRefundRequest struct {
	OrderDetail []*OrderDetailRefund    `json:"order_detail" binding:"required"`
	Reason      string                  `json:"reason" binding:"required"`
	RefundMedia []*multipart.FileHeader `form:"refund_media"`
}

type CreateRefundResponse struct {
	ID                    uuid.UUID `json:"id"`
	OrderID               uuid.UUID `json:"order_id"`
	Reason                string    `json:"reason"`
	RefundRequestDatetime time.Time `json:"refund_request_datetime"`
	Status                string    `json:"status"`
}

func NewCreateRefundResponse(source *models.Refund) (target *CreateRefundResponse) {
	target = &CreateRefundResponse{
		ID:                    source.ID,
		OrderID:               source.OrderID,
		Reason:                source.Reason,
		RefundRequestDatetime: source.RefundRequestDatetime,
		Status:                source.Status,
	}
	return
}
