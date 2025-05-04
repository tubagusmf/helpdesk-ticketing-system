package usecase

import (
	"context"
	"errors"
	"helpdesk-ticketing-system/internal/model"

	"github.com/sirupsen/logrus"
)

type ticketHistoryUsecase struct {
	ticketHistoryRepo model.ITicketHistoryRepository
}

func NewTicketHistoryUsecase(ticketHistoryRepo model.ITicketHistoryRepository) model.ITicketHistoryUsecase {
	return &ticketHistoryUsecase{ticketHistoryRepo: ticketHistoryRepo}
}

func (t *ticketHistoryUsecase) GetTicketID(ctx context.Context, id int64) (*model.TicketHistory, error) {
	log := logrus.WithFields(logrus.Fields{
		"id": id,
	})

	ticketHistory, err := t.ticketHistoryRepo.GetTicketID(ctx, id)
	if err != nil {
		log.Error("Failed to fetch ticket history by ID: ", err)
		return nil, err
	}

	if ticketHistory == nil {
		log.Error("Ticket history not found")
		return nil, errors.New("ticket history not found")
	}

	return ticketHistory, nil
}

func (t *ticketHistoryUsecase) GetStatus(ctx context.Context, status string) (*[]model.TicketHistory, error) {
	log := logrus.WithFields(logrus.Fields{
		"status": status,
	})

	ticketHistory, err := t.ticketHistoryRepo.GetStatus(ctx, status)
	if err != nil {
		log.Error("Failed to fetch ticket history by status: ", err)
		return nil, err
	}

	if ticketHistory == nil {
		log.Error("Ticket history not found")
		return nil, errors.New("ticket history not found")
	}

	return ticketHistory, nil
}

func (t *ticketHistoryUsecase) GetPriority(ctx context.Context, priority string) (*[]model.TicketHistory, error) {
	log := logrus.WithFields(logrus.Fields{
		"priority": priority,
	})

	ticketHistory, err := t.ticketHistoryRepo.GetPriority(ctx, priority)
	if err != nil {
		log.Error("Failed to fetch ticket history by priority: ", err)
		return nil, err
	}

	if ticketHistory == nil {
		log.Error("Ticket history not found")
		return nil, errors.New("ticket history not found")
	}

	return ticketHistory, nil
}

func (t *ticketHistoryUsecase) GetUserID(ctx context.Context, userID int64) (*[]model.TicketHistory, error) {
	log := logrus.WithFields(logrus.Fields{
		"user_id": userID,
	})

	ticketHistory, err := t.ticketHistoryRepo.GetUserID(ctx, userID)
	if err != nil {
		log.Error("Failed to fetch ticket history by user ID: ", err)
		return nil, err
	}

	if ticketHistory == nil {
		log.Error("Ticket history not found")
		return nil, errors.New("ticket history not found")
	}

	return ticketHistory, nil
}
