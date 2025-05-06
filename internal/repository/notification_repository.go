package repository

import (
	"context"
	"helpdesk-ticketing-system/internal/model"

	"gorm.io/gorm"
)

type NotificationRepo struct {
	db *gorm.DB
}

func NewNotificationRepo(db *gorm.DB) model.INotificationRepository {
	return &NotificationRepo{db: db}
}

func (n *NotificationRepo) Save(ctx context.Context, notification *model.Notification) error {
	err := n.db.WithContext(ctx).Create(notification).Error
	if err != nil {
		return err
	}

	return nil
}
