package repository

import (
	"context"
	"helpdesk-ticketing-system/internal/model"

	"gorm.io/gorm"
)

type AttachmentRepo struct {
	db *gorm.DB
}

func NewAttachmentRepo(db *gorm.DB) model.IAttachmentRepository {
	return &AttachmentRepo{db: db}
}

func (a *AttachmentRepo) FindAllByTicketID(ctx context.Context, ticketID int64) ([]*model.Attachment, error) {
	var attachments []*model.Attachment
	query := a.db.WithContext(ctx).Model(&model.Attachment{})

	err := query.Where("ticket_id = ?", ticketID).Order("uploaded_at ASC").Find(&attachments).Error

	return attachments, err
}

func (a *AttachmentRepo) Create(ctx context.Context, attachment model.Attachment) error {
	err := a.db.WithContext(ctx).Create(&attachment).Error
	if err != nil {
		return err
	}

	return nil
}
