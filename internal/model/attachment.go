package model

import (
	"context"
	"time"
)

type Attachment struct {
	ID         int64     `json:"id"`
	TicketID   int64     `json:"ticket_id"`
	FilePath   string    `json:"file_path"`
	UploadedAt time.Time `json:"uploaded_at"`
}

type AttachmentResponse struct {
	ID         int64     `json:"id"`
	TicketID   int64     `json:"ticket_id"`
	FilePath   string    `json:"file_path"`
	UploadedAt time.Time `json:"uploaded_at"`
}

type CreateAttachmentInput struct {
	TicketID int64  `json:"ticket_id" validate:"required"`
	FilePath string `json:"file_path" validate:"required"`
}

type IAttachmentRepository interface {
	FindAllByTicketID(ctx context.Context, ticketID int64) ([]*Attachment, error)
	Create(ctx context.Context, attachment Attachment) error
}

type IAttachmentUsecase interface {
	FindAllByTicketID(ctx context.Context, ticketID int64) ([]*Attachment, error)
	Create(ctx context.Context, in CreateAttachmentInput) error
}
