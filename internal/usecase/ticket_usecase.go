package usecase

import (
	"context"
	"errors"
	"helpdesk-ticketing-system/internal/helper"
	"helpdesk-ticketing-system/internal/model"
	"time"

	"github.com/sirupsen/logrus"
)

type TicketUsecase struct {
	ticketRepo  model.ITicketRepository
	userRepo    model.IUserRepository
	commentRepo model.ICommentRepository
}

func NewTicketUsecase(ticketRepo model.ITicketRepository, userRepo model.IUserRepository, commentRepo model.ICommentRepository) model.ITicketUsecase {
	return &TicketUsecase{
		ticketRepo:  ticketRepo,
		userRepo:    userRepo,
		commentRepo: commentRepo,
	}
}

func (t *TicketUsecase) FindAll(ctx context.Context, filter model.FindAllParam) ([]*model.TicketResponse, error) {
	log := logrus.WithFields(logrus.Fields{
		"filter": filter,
	})

	tickets, err := t.ticketRepo.FindAll(ctx, filter)
	if err != nil {
		log.Error("Failed to fetch tickets: ", err)
		return nil, err
	}

	var responses []*model.TicketResponse

	for _, ticket := range tickets {
		user, _ := t.userRepo.FindById(ctx, ticket.UserID)
		comments, _ := t.commentRepo.FindAllByTicketID(ctx, ticket.ID)

		var userRes *model.UserResponse
		if user != nil {
			userRes = &model.UserResponse{
				Name:  user.Name,
				Email: user.Email,
			}
		}

		var commentResList []*model.CommentResponse
		for _, comment := range comments {
			commentResList = append(commentResList, &model.CommentResponse{
				UserID:  comment.UserID,
				Content: comment.Content,
			})
		}

		response := &model.TicketResponse{
			ID:          ticket.ID,
			Title:       ticket.Title,
			Description: ticket.Description,
			Status:      ticket.Status,
			Priority:    ticket.Priority,
			AssignedTo:  ticket.AssignedTo,
			UserID:      ticket.UserID,
			User:        userRes,
			Comment:     commentResList,
			DueBy:       ticket.DueBy,
			CreatedAt:   ticket.CreatedAt,
			UpdatedAt:   ticket.UpdatedAt,
		}

		responses = append(responses, response)
	}

	return responses, nil
}

func (t *TicketUsecase) FindById(ctx context.Context, id int64) (*model.TicketResponse, error) {
	log := logrus.WithFields(logrus.Fields{
		"id": id,
	})

	ticket, err := t.ticketRepo.FindById(ctx, id)
	if err != nil {
		log.Error("Failed to fetch ticket by ID: ", err)
		return nil, err
	}

	if ticket == nil {
		log.Error("Ticket not found")
		return nil, errors.New("ticket not found")
	}

	user, _ := t.userRepo.FindById(ctx, ticket.UserID)
	comments, _ := t.commentRepo.FindAllByTicketID(ctx, ticket.ID)

	var userRes *model.UserResponse
	if user != nil {
		userRes = &model.UserResponse{
			Name:  user.Name,
			Email: user.Email,
		}
	}

	var commentResList []*model.CommentResponse
	for _, comment := range comments {
		commentResList = append(commentResList, &model.CommentResponse{
			UserID:  comment.UserID,
			Content: comment.Content,
		})
	}

	response := &model.TicketResponse{
		ID:          ticket.ID,
		Title:       ticket.Title,
		Description: ticket.Description,
		Status:      ticket.Status,
		Priority:    ticket.Priority,
		AssignedTo:  ticket.AssignedTo,
		UserID:      ticket.UserID,
		User:        userRes,
		Comment:     commentResList,
		DueBy:       ticket.DueBy,
		CreatedAt:   ticket.CreatedAt,
		UpdatedAt:   ticket.UpdatedAt,
	}

	return response, nil
}

func (t *TicketUsecase) Create(ctx context.Context, in model.CreateTicketInput) (*model.Ticket, error) {
	log := logrus.WithFields(logrus.Fields{
		"input": in,
	})

	err := helper.Validator.Struct(in)
	if err != nil {
		log.Error("Validation error: ", err)
		return &model.Ticket{}, err
	}

	userID, err := helper.GetUserID(ctx)
	if err != nil {
		log.Error("Failed to get user ID: ", err)
		return &model.Ticket{}, err
	}

	ticket := model.Ticket{
		Title:       in.Title,
		Description: in.Description,
		Status:      in.Status,
		Priority:    in.Priority,
		AssignedTo:  in.AssignedTo,
		UserID:      userID,
		DueBy:       helper.CalculateDueBy(in.Priority),
	}

	tickets, err := t.ticketRepo.Create(ctx, ticket)
	if err != nil {
		log.Error("Failed to create ticket: ", err)
		return &model.Ticket{}, err
	}

	return tickets, nil
}

func (t *TicketUsecase) Update(ctx context.Context, id int64, in model.UpdateTicketInput) (*model.Ticket, error) {
	log := logrus.WithFields(logrus.Fields{
		"id":    id,
		"input": in,
	})

	err := helper.Validator.Struct(in)
	if err != nil {
		log.Error("Validation error: ", err)
		return &model.Ticket{}, err
	}

	userID, err := helper.GetUserID(ctx)
	if err != nil {
		log.Error("Failed to get user ID: ", err)
		return &model.Ticket{}, err
	}

	exitingTicket, err := t.ticketRepo.FindById(ctx, id)
	if err != nil {
		log.Error("Failed to fetch ticket: ", err)
		return &model.Ticket{}, err
	}

	if exitingTicket == nil || (exitingTicket.DeletedAt != nil && !exitingTicket.DeletedAt.IsZero()) {
		log.Error("Ticket is deleted or does not exist")
		return &model.Ticket{}, errors.New("ticket is deleted or does not exist")
	}

	func(ticket *model.Ticket, input model.UpdateTicketInput) {
		ticket.Title = input.Title
		ticket.Description = input.Description
		ticket.Status = input.Status
		ticket.Priority = input.Priority
		ticket.AssignedTo = input.AssignedTo
		ticket.UserID = userID
		ticket.DueBy = helper.CalculateDueBy(input.Priority)
		ticket.UpdatedAt = time.Now()
	}(exitingTicket, in)

	tickets, err := t.ticketRepo.Update(ctx, *exitingTicket)
	if err != nil {
		log.Error("Failed to update ticket: ", err)
		return &model.Ticket{}, err
	}

	return tickets, nil
}

func (t *TicketUsecase) Delete(ctx context.Context, id int64) error {
	log := logrus.WithFields(logrus.Fields{
		"id": id,
	})

	ticket, err := t.ticketRepo.FindById(ctx, id)
	if err != nil {
		log.Error("Failed to fetch ticket: ", err)
		return err
	}

	if ticket == nil {
		log.Error("Ticket not found")
		return errors.New("ticket not found")
	}

	if ticket.DeletedAt != nil && !ticket.DeletedAt.IsZero() {
		log.Error("Ticket is already deleted")
		return errors.New("ticket is already deleted")
	}

	err = t.ticketRepo.Delete(ctx, id)
	if err != nil {
		log.Error("Failed to delete ticket: ", err)
		return err
	}

	log.Info("Successfully deleted ticket with ID: ", id)
	return nil
}
