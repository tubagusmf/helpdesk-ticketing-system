package model

import (
	"context"
	"time"
)

type ITicketRepository interface {
	FindAll(ctx context.Context, filter FindAllParam) ([]*TicketResponse, error)
	FindById(ctx context.Context, id int64) (*Ticket, error)
	Create(ctx context.Context, ticket Ticket) (*Ticket, error)
	Update(ctx context.Context, ticket Ticket) (*Ticket, error)
	Delete(ctx context.Context, id int64) error
}

type ITicketUsecase interface {
	FindAll(ctx context.Context, filter FindAllParam) ([]*TicketResponse, error)
	FindById(ctx context.Context, id int64) (*TicketResponse, error)
	Create(ctx context.Context, in CreateTicketInput) (*Ticket, error)
	Update(ctx context.Context, id int64, in UpdateTicketInput) (*Ticket, error)
	Delete(ctx context.Context, id int64) error
}

type Ticket struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Priority    string     `json:"priority"`
	AssignedTo  int64      `json:"assigned_to"`
	UserID      int64      `json:"user_id"`
	DueBy       *time.Time `json:"due_by,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"-"`
}

type TicketResponse struct {
	ID          int64                 `json:"id"`
	Title       string                `json:"title"`
	Description string                `json:"description"`
	Status      string                `json:"status"`
	Priority    string                `json:"priority"`
	AssignedTo  int64                 `json:"assigned_to"`
	UserID      int64                 `json:"user_id"`
	User        *UserResponse         `json:"user,omitempty"`
	Comment     []*CommentResponse    `json:"comment,omitempty"`
	Attachment  []*AttachmentResponse `json:"attachment,omitempty"`
	DueBy       *time.Time            `json:"due_by,omitempty"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
}

type FindAllParam struct {
	Limit int64 `json:"limit"`
	Page  int64 `json:"page"`
}

type CreateTicketInput struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	Status      string `json:"status" validate:"required"`
	Priority    string `json:"priority" validate:"required"`
	AssignedTo  int64  `json:"assigned_to" validate:"required"`
}

type UpdateTicketInput struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	Status      string `json:"status" validate:"required"`
	Priority    string `json:"priority" validate:"required"`
	AssignedTo  int64  `json:"assigned_to" validate:"required"`
}
