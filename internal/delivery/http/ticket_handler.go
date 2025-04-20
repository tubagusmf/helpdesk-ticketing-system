package http

import (
	"helpdesk-ticketing-system/internal/model"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type TicketHandler struct {
	ticketUsecase model.ITicketUsecase
}

func NewTicketHandler(e *echo.Echo, ticketUsecase model.ITicketUsecase) {
	handler := &TicketHandler{ticketUsecase: ticketUsecase}

	routeUrl := e.Group("v1/ticket")
	routeUrl.GET("", handler.FindAll, AuthMiddleware)
	routeUrl.GET("/:id", handler.FindById, AuthMiddleware)
	routeUrl.POST("/create", handler.Create, AuthMiddleware)
	routeUrl.PUT("/update/:id", handler.Update, AuthMiddleware)
	routeUrl.DELETE("/delete/:id", handler.Delete, AuthMiddleware)
}

func (h *TicketHandler) FindAll(c echo.Context) error {
	tickets, err := h.ticketUsecase.FindAll(c.Request().Context(), model.FindAllParam{})

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, Response{
		Status: http.StatusOK,
		Data:   tickets,
	})
}

func (h *TicketHandler) FindById(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ticket ID format")
	}

	ticket, err := h.ticketUsecase.FindById(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Ticket not found")
	}

	return c.JSON(http.StatusOK, Response{
		Status: http.StatusOK,
		Data:   ticket,
	})
}

func (h *TicketHandler) Create(c echo.Context) error {
	var body model.CreateTicketInput
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	ticket, err := h.ticketUsecase.Create(c.Request().Context(), body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create ticket")
	}

	return c.JSON(http.StatusCreated, Response{
		Status: http.StatusCreated,
		Data:   ticket,
	})
}

func (h *TicketHandler) Update(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ticket ID format")
	}

	var body model.UpdateTicketInput
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	ticket, err := h.ticketUsecase.Update(c.Request().Context(), id, body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update ticket")
	}

	return c.JSON(http.StatusOK, Response{
		Status:  http.StatusOK,
		Message: "Ticket updated successfully",
		Data:    ticket,
	})
}

func (h *TicketHandler) Delete(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ticket ID format")
	}

	err = h.ticketUsecase.Delete(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete ticket")
	}

	return c.JSON(http.StatusOK, Response{
		Status:  http.StatusOK,
		Message: "Ticket deleted successfully",
	})
}
