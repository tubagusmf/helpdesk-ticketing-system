package usecase

import (
	"context"
	"errors"
	"fmt"
	"helpdesk-ticketing-system/internal/helper"
	"helpdesk-ticketing-system/internal/model"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type TicketUsecase struct {
	ticketRepo          model.ITicketRepository
	userRepo            model.IUserRepository
	commentRepo         model.ICommentRepository
	attachmentRepo      model.IAttachmentRepository
	ticketHistoryRepo   model.ITicketHistoryRepository
	notificationUsecase model.INotificationUsecase
	rmq                 *amqp.Channel
}

func NewTicketUsecase(
	ticketRepo model.ITicketRepository,
	userRepo model.IUserRepository,
	commentRepo model.ICommentRepository,
	attachmentRepo model.IAttachmentRepository,
	ticketHistoryRepo model.ITicketHistoryRepository,
	notificationUsecase model.INotificationUsecase,
	rmq *amqp.Channel,
) model.ITicketUsecase {
	return &TicketUsecase{
		ticketRepo:          ticketRepo,
		userRepo:            userRepo,
		commentRepo:         commentRepo,
		attachmentRepo:      attachmentRepo,
		ticketHistoryRepo:   ticketHistoryRepo,
		notificationUsecase: notificationUsecase,
		rmq:                 rmq,
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
		attachments, _ := t.attachmentRepo.FindAllByTicketID(ctx, ticket.ID)

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

		var attachmentResList []*model.AttachmentResponseForTicket
		for _, attachment := range attachments {
			attachmentResList = append(attachmentResList, &model.AttachmentResponseForTicket{
				FilePath:   attachment.FilePath,
				UploadedAt: attachment.UploadedAt,
			})
		}

		penalty := false
		var overdueBy string

		if (ticket.Status == "open" || ticket.Status == "in_progress") && time.Now().After(*ticket.DueBy) {
			penalty = true
			overdueDuration := time.Since(*ticket.DueBy)
			overdueBy = helper.FormatDuration(overdueDuration)
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
			Attachment:  attachmentResList,
			DueBy:       ticket.DueBy,
			CreatedAt:   ticket.CreatedAt,
			UpdatedAt:   ticket.UpdatedAt,
			Penalty:     penalty,
			Overdueby:   overdueBy,
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
	attachments, _ := t.attachmentRepo.FindAllByTicketID(ctx, ticket.ID)

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

	var attachmentResList []*model.AttachmentResponseForTicket
	for _, attachment := range attachments {
		attachmentResList = append(attachmentResList, &model.AttachmentResponseForTicket{
			FilePath:   attachment.FilePath,
			UploadedAt: attachment.UploadedAt,
		})
	}

	penalty := false
	var overdueBy string

	if (ticket.Status == "open" || ticket.Status == "in_progress") && time.Now().After(*ticket.DueBy) {
		penalty = true
		overdueDuration := time.Since(*ticket.DueBy)
		overdueBy = helper.FormatDuration(overdueDuration)
	}

	response := &model.TicketResponse{
		ID:          ticket.ID,
		Title:       ticket.Title,
		Description: ticket.Description,
		Status:      ticket.Status,
		Priority:    ticket.Priority,
		AssignedTo:  ticket.AssignedTo,
		UserID:      ticket.UserID,
		Attachment:  attachmentResList,
		User:        userRes,
		Comment:     commentResList,
		DueBy:       ticket.DueBy,
		CreatedAt:   ticket.CreatedAt,
		UpdatedAt:   ticket.UpdatedAt,
		Penalty:     penalty,
		Overdueby:   overdueBy,
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

	assignedUser, err := t.userRepo.FindById(ctx, in.AssignedTo)
	if err != nil {
		log.Error("Failed to fetch assigned user: ", err)
		return nil, fmt.Errorf("failed to fetch assigned user")
	}

	notification := model.Notification{
		UserID:    assignedUser.ID,
		Email:     assignedUser.Email,
		Subject:   tickets.Title,
		Message:   tickets.Description,
		Status:    "pending",
		TicketID:  tickets.ID,
		CreatedAt: time.Now(),
	}

	err = t.notificationUsecase.SendNotification(ctx, &notification)
	if err != nil {
		log.Warn("Failed to send notification: ", err)
		return nil, err
	}

	ticketHistory := model.TicketHistory{
		TicketID:  tickets.ID,
		UserID:    userID,
		Status:    tickets.Status,
		Priority:  tickets.Priority,
		ChangedAt: time.Now(),
	}

	err = t.ticketHistoryRepo.Create(ctx, ticketHistory)
	if err != nil {
		return nil, err
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

	ticketHistory := model.TicketHistory{
		TicketID:  tickets.ID,
		UserID:    tickets.UserID,
		Status:    tickets.Status,
		Priority:  tickets.Priority,
		ChangedAt: time.Now(),
	}
	err = t.ticketHistoryRepo.Create(ctx, ticketHistory)
	if err != nil {
		return nil, err
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
