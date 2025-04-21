package usecase

import (
	"context"
	"errors"
	"helpdesk-ticketing-system/internal/helper"
	"helpdesk-ticketing-system/internal/model"

	"github.com/sirupsen/logrus"
)

type CommentUsecase struct {
	commentRepo model.ICommentRepository
}

func NewCommentUsecase(commentRepo model.ICommentRepository) model.ICommentUsecase {
	return &CommentUsecase{commentRepo: commentRepo}
}

func (c *CommentUsecase) FindAll(ctx context.Context, comment model.Comment) ([]*model.Comment, error) {
	log := logrus.WithFields(logrus.Fields{
		"comment": comment,
	})

	comments, err := c.commentRepo.FindAll(ctx, comment)
	if err != nil {
		log.Error("Failed to fetch comments: ", err)
		return nil, err
	}

	return comments, nil
}

func (c *CommentUsecase) FindById(ctx context.Context, id int64) (*model.Comment, error) {
	log := logrus.WithFields(logrus.Fields{
		"id": id,
	})

	comment, err := c.commentRepo.FindById(ctx, id)
	if err != nil {
		log.Error("Failed to fetch comment by ID: ", err)
		return nil, err
	}

	if comment == nil {
		log.Error("Comment not found")
		return nil, errors.New("comment not found")
	}

	return comment, nil
}

func (c *CommentUsecase) Create(ctx context.Context, in model.CreateCommentInput) (*model.Comment, error) {
	log := logrus.WithFields(logrus.Fields{
		"input": in,
	})

	err := helper.Validator.Struct(in)
	if err != nil {
		log.Error("Validation error: ", err)
		return &model.Comment{}, err
	}

	userID, err := helper.GetUserID(ctx)
	if err != nil {
		log.Error("Failed to get user ID: ", err)
		return &model.Comment{}, err
	}

	comment := model.Comment{
		UserID:   userID,
		TicketID: in.TicketId,
		Content:  in.Content,
	}

	comments, err := c.commentRepo.Create(ctx, comment)
	if err != nil {
		log.Error("Failed to create comment: ", err)
		return &model.Comment{}, err
	}

	return comments, nil
}

func (c *CommentUsecase) Update(ctx context.Context, id int64, in model.UpdateCommentInput) (*model.Comment, error) {
	log := logrus.WithFields(logrus.Fields{
		"input": in,
	})

	err := helper.Validator.Struct(in)
	if err != nil {
		log.Error("Validation error: ", err)
		return &model.Comment{}, err
	}

	userID, err := helper.GetUserID(ctx)
	if err != nil {
		log.Error("Failed to get user ID: ", err)
		return &model.Comment{}, err
	}

	exitingComment, err := c.commentRepo.FindById(ctx, id)
	if err != nil {
		log.Error("Failed to fetch comment: ", err)
		return &model.Comment{}, err
	}

	if exitingComment == nil || (exitingComment.DeletedAt != nil && !exitingComment.DeletedAt.IsZero()) {
		log.Error("Comment is deleted or does not exist")
		return &model.Comment{}, errors.New("comment is deleted or does not exist")
	}

	if exitingComment.UserID != userID {
		log.Error("You are not authorized to update this comment")
		return &model.Comment{}, errors.New("you are not authorized to update this comment")
	}

	comments, err := c.commentRepo.Update(ctx, model.Comment{
		ID:        id,
		UserID:    userID,
		TicketID:  in.TicketId,
		Content:   in.Content,
		DeletedAt: nil,
	})
	if err != nil {
		log.Error("Failed to update comment: ", err)
		return &model.Comment{}, err
	}

	return comments, nil
}

func (c *CommentUsecase) Delete(ctx context.Context, id int64) error {
	log := logrus.WithFields(logrus.Fields{
		"id": id,
	})

	comment, err := c.commentRepo.FindById(ctx, id)
	if err != nil {
		log.Error("Failed to fetch comment by ID: ", err)
		return err
	}

	if comment == nil {
		log.Error("Comment not found")
		return errors.New("comment not found")
	}

	if comment.DeletedAt != nil && !comment.DeletedAt.IsZero() {
		log.Error("Comment is already deleted")
		return errors.New("comment is already deleted")
	}

	err = c.commentRepo.Delete(ctx, id)
	if err != nil {
		log.Error("Failed to delete comment: ", err)
		return err
	}

	log.Info("Successfully deleted comment with ID: ", id)
	return nil
}
