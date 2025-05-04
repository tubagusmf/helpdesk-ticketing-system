package repository

import (
	"context"
	"helpdesk-ticketing-system/internal/model"

	"gorm.io/gorm"
)

type TicketHistoryRepo struct {
	db *gorm.DB
}

func NewTicketHistoryRepo(db *gorm.DB) model.ITicketHistoryRepository {
	return &TicketHistoryRepo{db: db}
}

func (t *TicketHistoryRepo) GetTicketID(ctx context.Context, id int64) (*model.TicketHistory, error) {
	var history model.TicketHistory

	err := t.db.WithContext(ctx).
		Where("ticket_id = ?", id).
		Order("changed_at DESC").
		First(&history).Error
	if err != nil {
		return nil, err
	}

	return &history, err
}

func (t *TicketHistoryRepo) GetStatus(ctx context.Context, status string) (*[]model.TicketHistory, error) {
	var histories []model.TicketHistory

	err := t.db.WithContext(ctx).
		Where("status = ?", status).
		Order("changed_at DESC").
		Find(&histories).Error

	return &histories, err
}

func (t *TicketHistoryRepo) GetPriority(ctx context.Context, priority string) (*[]model.TicketHistory, error) {
	var histories []model.TicketHistory

	err := t.db.WithContext(ctx).
		Where("priority = ?", priority).
		Order("changed_at DESC").
		Find(&histories).Error

	return &histories, err
}

func (t *TicketHistoryRepo) GetUserID(ctx context.Context, userID int64) (*[]model.TicketHistory, error) {
	var histories []model.TicketHistory

	err := t.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("changed_at DESC").
		Find(&histories).Error

	return &histories, err
}

func (t *TicketHistoryRepo) Create(ctx context.Context, ticketHistory model.TicketHistory) error {
	err := t.db.WithContext(ctx).Create(&ticketHistory).Error
	if err != nil {
		return err
	}

	return nil
}
