package model

import (
	"context"
	"time"
)

type TicketHistory struct {
	ID        int64     `json:"id"`
	TicketID  int64     `json:"ticket_id"`
	UserID    int64     `json:"user_id"`
	Status    string    `json:"status"`
	Priority  string    `json:"priority"`
	ChangedAt time.Time `json:"changed_at"`
}

type ITicketHistoryRepository interface {
	GetTicketID(ctx context.Context, id int64) (*TicketHistory, error)
	GetStatus(ctx context.Context, status string) (*[]TicketHistory, error)
	GetPriority(ctx context.Context, priority string) (*[]TicketHistory, error)
	GetUserID(ctx context.Context, userID int64) (*[]TicketHistory, error)
	Create(ctx context.Context, ticketHistory TicketHistory) error
}

type ITicketHistoryUsecase interface {
	GetTicketID(ctx context.Context, id int64) (*TicketHistory, error)
	GetStatus(ctx context.Context, status string) (*[]TicketHistory, error)
	GetPriority(ctx context.Context, priority string) (*[]TicketHistory, error)
	GetUserID(ctx context.Context, userID int64) (*[]TicketHistory, error)
}
