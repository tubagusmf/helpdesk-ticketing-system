package model

import (
	"context"
	"time"
)

type Comment struct {
	ID        int64      `json:"id"`
	UserID    int64      `json:"user_id"`
	TicketID  int64      `json:"ticket_id"`
	Content   string     `json:"content"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
}

type CreateCommentInput struct {
	TicketId int64  `json:"ticket_id" validate:"required"`
	Content  string `json:"content" validate:"required"`
}

type UpdateCommentInput struct {
	TicketId int64  `json:"ticket_id" validate:"required"`
	Content  string `json:"content" validate:"required"`
}

type ICommentRepository interface {
	FindAll(ctx context.Context, comment Comment) ([]*Comment, error)
	FindById(ctx context.Context, id int64) (*Comment, error)
	Create(ctx context.Context, comment Comment) (*Comment, error)
	Update(ctx context.Context, comment Comment) (*Comment, error)
	Delete(ctx context.Context, id int64) error
}

type ICommentUsecase interface {
	FindAll(ctx context.Context, comment Comment) ([]*Comment, error)
	FindById(ctx context.Context, id int64) (*Comment, error)
	Create(ctx context.Context, in CreateCommentInput) (*Comment, error)
	Update(ctx context.Context, id int64, in UpdateCommentInput) (*Comment, error)
	Delete(ctx context.Context, id int64) error
}
