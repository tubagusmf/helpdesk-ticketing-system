package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"helpdesk-ticketing-system/internal/model"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const (
	cacheKeyAll  = "tickets:all"
	cacheKeyByID = "ticket:%d"
)

type TaskRepo struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewTicketRepo(db *gorm.DB, rdb *redis.Client) model.ITicketRepository {
	return &TaskRepo{
		db:  db,
		rdb: rdb,
	}
}

func (t *TaskRepo) FindAll(ctx context.Context, filter model.FindAllParam) ([]*model.TicketResponse, error) {
	cached, err := t.rdb.Get(ctx, cacheKeyAll).Result()
	if err == nil {
		var tickets []*model.TicketResponse
		if err := json.Unmarshal([]byte(cached), &tickets); err == nil {
			return tickets, nil
		}
	}

	var tickets []*model.TicketResponse
	query := t.db.WithContext(ctx).Model(&model.Ticket{})

	if filter.Limit > 0 {
		query = query.Limit(int(filter.Limit))
	}
	if filter.Page > 0 {
		offset := int((filter.Page - 1) * filter.Limit)
		query = query.Offset(offset)
	}

	err = query.Where("deleted_at IS NULL").Find(&tickets).Error
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(tickets)
	if err == nil {
		t.rdb.Set(ctx, cacheKeyAll, data, time.Minute*5)
	}

	return tickets, nil
}

func (t *TaskRepo) FindById(ctx context.Context, id int64) (*model.Ticket, error) {
	cacheKey := fmt.Sprintf(cacheKeyByID, id)

	cached, err := t.rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		var ticket model.Ticket
		if err := json.Unmarshal([]byte(cached), &ticket); err == nil {
			return &ticket, nil
		}
	}

	var ticket model.Ticket
	err = t.db.WithContext(ctx).Where("deleted_at IS NULL").First(&ticket, id).Preload("Comments").Error
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(ticket)
	if err == nil {
		t.rdb.Set(ctx, cacheKey, data, time.Minute*5)
	}

	return &ticket, nil
}

func (t *TaskRepo) Create(ctx context.Context, ticket model.Ticket) (*model.Ticket, error) {
	err := t.db.WithContext(ctx).Create(&ticket).Error
	if err != nil {
		return nil, err
	}

	t.rdb.Del(ctx, cacheKeyAll)

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

	t.rdb.Del(ctx, fmt.Sprintf(cacheKeyByID, ticket.ID))
	t.rdb.Del(ctx, cacheKeyAll)

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

	t.rdb.Del(ctx, fmt.Sprintf(cacheKeyByID, id))
	t.rdb.Del(ctx, cacheKeyAll)

	return nil
}
