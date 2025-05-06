package http

import (
	"helpdesk-ticketing-system/internal/model"
	"net/http"

	"github.com/labstack/echo/v4"
)

type NotificationHandler struct {
	notificationUsecase model.INotificationUsecase
}

func NewNotificationHandler(e *echo.Echo, notificationUsecase model.INotificationUsecase) {
	handler := &NotificationHandler{notificationUsecase: notificationUsecase}

	routeUrl := e.Group("v1/notification")
	routeUrl.POST("/send", handler.Send, AuthMiddleware)
}

func (n *NotificationHandler) Send(c echo.Context) error {
	var body model.Notification
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	err := n.notificationUsecase.SendNotification(c.Request().Context(), &body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to send notification")
	}

	return c.JSON(http.StatusOK, Response{
		Status:  http.StatusOK,
		Message: "Notification sent successfully",
	})
}
