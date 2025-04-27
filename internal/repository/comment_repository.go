package repository

import (
	"context"
	"helpdesk-ticketing-system/internal/model"
	"time"

	"gorm.io/gorm"
)

type CommentRepo struct {
	db *gorm.DB
}

func NewCommentRepo(db *gorm.DB) model.ICommentRepository {
	return &CommentRepo{db: db}
}

func (c *CommentRepo) FindAll(ctx context.Context, comment model.Comment) ([]*model.Comment, error) {
	var comments []*model.Comment
	query := c.db.WithContext(ctx).Model(&model.Comment{})

	err := query.Where("deleted_at IS NULL").Find(&comments).Error
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (c *CommentRepo) FindById(ctx context.Context, id int64) (*model.Comment, error) {
	var comment model.Comment

	err := c.db.WithContext(ctx).Where("deleted_at IS NULL").First(&comment, id).Error
	if err != nil {
		return nil, err
	}

	return &comment, nil
}

func (c *CommentRepo) FindAllByTicketID(ctx context.Context, ticketID int64) ([]*model.Comment, error) {
	var comments []*model.Comment

	err := c.db.WithContext(ctx).Where("ticket_id = ?", ticketID).Order("created_at ASC").Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (c *CommentRepo) Create(ctx context.Context, comment model.Comment) (*model.Comment, error) {
	err := c.db.WithContext(ctx).Create(&comment).Error
	if err != nil {
		return nil, err
	}

	return &comment, nil
}

func (c *CommentRepo) Update(ctx context.Context, comment model.Comment) (*model.Comment, error) {
	err := c.db.WithContext(ctx).
		Model(&model.Comment{}).
		Where("id = ?", comment.ID).
		Updates(&comment).Error

	if err != nil {
		return nil, err
	}

	return &comment, nil
}

func (c *CommentRepo) Delete(ctx context.Context, id int64) error {
	err := c.db.WithContext(ctx).
		Model(&model.Comment{}).
		Where("id = ?", id).
		Update("deleted_at", time.Now()).Error

	if err != nil {
		return err
	}

	return nil
}
