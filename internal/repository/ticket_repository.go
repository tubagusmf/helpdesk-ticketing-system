package repository

import (
	"context"
	"helpdesk-ticketing-system/internal/model"
	"time"

	"gorm.io/gorm"
)

type TaskRepo struct {
	db *gorm.DB
}

func NewTicketRepo(db *gorm.DB) model.ITicketRepository {
	return &TaskRepo{db: db}
}

func (t *TaskRepo) FindAll(ctx context.Context, filter model.FindAllParam) ([]*model.Ticket, error) {
	var ticket []*model.Ticket
	query := t.db.WithContext(ctx).Model(&model.Ticket{})

	if filter.Limit > 0 {
		query = query.Limit(int(filter.Limit))
	}
	if filter.Page > 0 {
		offset := int((filter.Page - 1) * filter.Limit)
		query = query.Offset(offset)
	}

	err := query.Where("deleted_at IS NULL").Find(&ticket).Error
	if err != nil {
		return nil, err
	}

	return ticket, nil
}

func (t *TaskRepo) FindById(ctx context.Context, id int64) (*model.Ticket, error) {
	var ticket model.Ticket
	err := t.db.WithContext(ctx).Where("deleted_at IS NULL").First(&ticket, id).Preload("Comments").Error

	if err != nil {
		return nil, err
	}

	return &ticket, nil
}

func (t *TaskRepo) Create(ctx context.Context, ticket model.Ticket) (*model.Ticket, error) {
	err := t.db.WithContext(ctx).Create(&ticket).Error
	if err != nil {
		return nil, err
	}

	return &ticket, nil
}

func (t *TaskRepo) Update(ctx context.Context, ticket model.Ticket) (*model.Ticket, error) {
	err := t.db.WithContext(ctx).
		Model(&model.Ticket{}).
		Where("id = ?", ticket.ID).
		Updates(&ticket).Error

	if err != nil {
		return nil, err
	}

	return &ticket, nil
}

func (t *TaskRepo) Delete(ctx context.Context, id int64) error {
	err := t.db.WithContext(ctx).
		Model(&model.Ticket{}).
		Where("id = ?", id).
		Update("deleted_at", time.Now()).Error

	if err != nil {
		return err
	}

	return nil
}
