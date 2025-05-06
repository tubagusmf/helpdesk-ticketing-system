package console

import (
	"helpdesk-ticketing-system/internal/config"
	"helpdesk-ticketing-system/internal/repository"
	"helpdesk-ticketing-system/internal/usecase"
	"log"
	"net/http"
	"sync"

	"helpdesk-ticketing-system/database"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"helpdesk-ticketing-system/internal/worker"

	handlerHttp "helpdesk-ticketing-system/internal/delivery/http"
)

func init() {
	rootCmd.AddCommand(serverCMD)
}

var serverCMD = &cobra.Command{
	Use:   "httpsrv",
	Short: "Start HTTP server",
	Long:  "Start the HTTP server to handle incoming requests for the to-do list application.",
	Run:   httpServer,
}

func httpServer(cmd *cobra.Command, args []string) {
	config.LoadWithViper()
	config.LoadWithGetenv()

	postgresDB := database.NewPostgres()
	sqlDB, err := postgresDB.DB()
	if err != nil {
		log.Fatalf("Failed to get SQL DB from Gorm: %v", err)
	}

	defer sqlDB.Close()

	rmqChannel, err := config.InitRabbitMQ()
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ channel: %v", err)
	}
	defer rmqChannel.Close()

	worker.StartEmailWorker(rmqChannel)

	userRepo := repository.NewUserRepo(postgresDB)
	userUsecase := usecase.NewUserUsecase(userRepo)
	commentRepo := repository.NewCommentRepo(postgresDB)
	commentUsecase := usecase.NewCommentUsecase(commentRepo)
	attachmentRepo := repository.NewAttachmentRepo(postgresDB)
	attachmentUsecase := usecase.NewAttachmentUsecase(attachmentRepo)
	ticketHistoryRepo := repository.NewTicketHistoryRepo(postgresDB)
	ticketHistoryUsecase := usecase.NewTicketHistoryUsecase(ticketHistoryRepo)
	notificationRepo := repository.NewNotificationRepo(postgresDB)
	notificationUsecase := usecase.NewNotificationUsecase(notificationRepo, rmqChannel)
	ticketRepo := repository.NewTicketRepo(postgresDB)
	ticketUsecase := usecase.NewTicketUsecase(
		ticketRepo,
		userRepo,
		commentRepo,
		attachmentRepo,
		ticketHistoryRepo, notificationUsecase,
		rmqChannel,
	)

	e := echo.New()

	handlerHttp.NewUserHandler(e, userUsecase)
	handlerHttp.NewTicketHandler(e, ticketUsecase)
	handlerHttp.NewCommentHandler(e, commentUsecase)
	handlerHttp.NewAttachmentHandler(e, attachmentUsecase)
	handlerHttp.NewTicketHistoryHandler(e, ticketHistoryUsecase)
	handlerHttp.NewNotificationHandler(e, notificationUsecase)

	var wg sync.WaitGroup
	errCh := make(chan error, 2)

	wg.Add(1)
	go func() {
		defer wg.Done()
		worker.StartEmailWorker(rmqChannel)

		select {}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		errCh <- e.Start(":3000")
	}()

	go func() {
		wg.Wait()
		close(errCh)
	}()

	if err := <-errCh; err != nil {
		if err != http.ErrServerClosed {
			logrus.Errorf("HTTP server error: %v", err)
		}
	}
}
