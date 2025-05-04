package http

import (
	"net/http"
	"strconv"

	"helpdesk-ticketing-system/internal/model"

	"github.com/labstack/echo/v4"
)

type TicketHistoryHandler struct {
	ticketHistoryUsecase model.ITicketHistoryUsecase
}

func NewTicketHistoryHandler(e *echo.Echo, ticketHistoryUsecase model.ITicketHistoryUsecase) {
	handler := &TicketHistoryHandler{ticketHistoryUsecase: ticketHistoryUsecase}

	routeUrl := e.Group("v1/ticket/history")
	routeUrl.GET("/id/:id", handler.GetByTicketID, AuthMiddleware)
	routeUrl.GET("/status/:status", handler.GetByStatus, AuthMiddleware)
	routeUrl.GET("/priority/:priority", handler.GetByPriority, AuthMiddleware)
	routeUrl.GET("/user/:id", handler.GetByUserID, AuthMiddleware)
}

func (t *TicketHistoryHandler) GetByTicketID(ctx echo.Context) error {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ticket ID")
	}

	history, err := t.ticketHistoryUsecase.GetTicketID(ctx.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Ticket not found")
	}

	return ctx.JSON(http.StatusOK, Response{
		Status: http.StatusOK,
		Data:   history,
	})
}

func (t *TicketHistoryHandler) GetByStatus(c echo.Context) error {
	status := c.Param("status")

	histories, err := t.ticketHistoryUsecase.GetStatus(c.Request().Context(), status)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Ticket not found")
	}

	return c.JSON(http.StatusOK, Response{
		Status: http.StatusOK,
		Data:   histories,
	})
}

func (h *TicketHistoryHandler) GetByPriority(c echo.Context) error {
	priority := c.Param("priority")

	histories, err := h.ticketHistoryUsecase.GetPriority(c.Request().Context(), priority)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Ticket not found")
	}

	return c.JSON(http.StatusOK, Response{
		Status: http.StatusOK,
		Data:   histories,
	})
}

func (t *TicketHistoryHandler) GetByUserID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	histories, err := t.ticketHistoryUsecase.GetUserID(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Ticket not found")
	}

	return c.JSON(http.StatusOK, Response{
		Status: http.StatusOK,
		Data:   histories,
	})
}
