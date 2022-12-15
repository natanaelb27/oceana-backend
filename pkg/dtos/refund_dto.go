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
	OrderDetail []*OrderDetailRefund  `json:"order_detail" binding:"required"`
	Reason      string                `json:"reason" binding:"required"`
	Photo1      *multipart.FileHeader `form:"photo_1"`
	Photo2      *multipart.FileHeader `form:"photo_2"`
	Photo3      *multipart.FileHeader `form:"photo_3"`
	Photo4      *multipart.FileHeader `form:"photo_4"`
	Photo5      *multipart.FileHeader `form:"photo_5"`
	Video       *multipart.FileHeader `form:"video"`
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
