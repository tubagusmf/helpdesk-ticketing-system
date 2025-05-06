package model

import (
	"context"
	"time"
)

type Notification struct {
	ID        int64     `json:"id"`
	TicketID  int64     `json:"ticket_id"`
	UserID    int64     `json:"user_id"`
	Email     string    `json:"email"`
	Subject   string    `json:"subject"`
	Message   string    `json:"message"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type INotificationRepository interface {
	Save(ctx context.Context, notification *Notification) error
}

type INotificationUsecase interface {
	SendNotification(ctx context.Context, notification *Notification) error
}
