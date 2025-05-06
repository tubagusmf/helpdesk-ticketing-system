package worker

import (
	"encoding/json"
	"helpdesk-ticketing-system/internal/model"
	"log"
	"net/smtp"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func StartEmailWorker(ch *amqp.Channel) {
	msgs, err := ch.Consume(
		"emailQueue",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println("Failed to register a consumer:", err)
		return
	}

	go func() {
		for d := range msgs {
			var notif model.Notification
			json.Unmarshal(d.Body, &notif)

			// call your email sending function
			SendEmail(notif.Email, notif.Subject, notif.Message)

			log.Println("Email sent to:", notif.Email)
		}
	}()
}

func SendEmail(to string, subject string, message string) {
	from := os.Getenv("EMAIL_FROM")
	password := os.Getenv("EMAIL_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	auth := smtp.PlainAuth("", from, password, smtpHost)

	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/plain; charset=\"utf-8\"\r\n" +
		"\r\n" +
		message + "\r\n")

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
	if err != nil {
		log.Println("Failed to send email:", err)
		return
	}

	log.Println("Email successfully sent to", to)
}
