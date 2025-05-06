package usecase

import (
	"context"
	"encoding/json"
	"helpdesk-ticketing-system/internal/model"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type NotificationUsecase struct {
	notificationRepo model.INotificationRepository
	rmq              *amqp.Channel
}

func NewNotificationUsecase(notificationRepo model.INotificationRepository, rmq *amqp.Channel) model.INotificationUsecase {
	return &NotificationUsecase{
		notificationRepo: notificationRepo,
		rmq:              rmq,
	}
}

func (n *NotificationUsecase) SendNotification(ctx context.Context, notification *model.Notification) error {
	log := logrus.WithFields(logrus.Fields{
		"notification": notification,
	})

	body, err := json.Marshal(notification)
	if err != nil {
		log.Error("Failed to marshal notification: ", err)
		return err
	}

	err = n.rmq.Publish(
		"notification", // if empty, use default exchange
		"emailQueue",   // if empty, use default routing key
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Error("Failed to publish notification: ", err)
		return err
	}

	err = n.notificationRepo.Save(ctx, notification)
	if err != nil {
		log.Error("Failed to save notification: ", err)
		return err
	}

	return nil
}
