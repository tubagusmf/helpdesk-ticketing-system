package usecase

import (
	"context"
	"helpdesk-ticketing-system/internal/helper"
	"helpdesk-ticketing-system/internal/model"
	"time"

	"github.com/sirupsen/logrus"
)

type AttachmentUsecase struct {
	attachmentRepo model.IAttachmentRepository
}

func NewAttachmentUsecase(attachmentRepo model.IAttachmentRepository) model.IAttachmentUsecase {
	return &AttachmentUsecase{
		attachmentRepo: attachmentRepo,
	}
}

func (a *AttachmentUsecase) FindAllByTicketID(ctx context.Context, ticketID int64) ([]*model.Attachment, error) {
	log := logrus.WithFields(logrus.Fields{
		"ticketID": ticketID,
	})

	attachments, err := a.attachmentRepo.FindAllByTicketID(ctx, ticketID)
	if err != nil {
		log.Error("Failed to fetch attachments: ", err)
		return nil, err
	}

	return attachments, nil
}

func (a *AttachmentUsecase) Create(ctx context.Context, in model.CreateAttachmentInput) error {
	log := logrus.WithFields(logrus.Fields{
		"input": in,
	})

	err := helper.Validator.Struct(in)
	if err != nil {
		log.Error("Validation error: ", err)
		return err
	}

	attachment := model.Attachment{
		TicketID:   in.TicketID,
		FilePath:   in.FilePath,
		UploadedAt: time.Now(),
	}

	err = a.attachmentRepo.Create(ctx, attachment)
	if err != nil {
		log.Error("Failed to create attachment: ", err)
		return err
	}

	return nil
}
