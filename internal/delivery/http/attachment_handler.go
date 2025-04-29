package http

import (
	"fmt"
	"helpdesk-ticketing-system/internal/model"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type AttachmentHandler struct {
	attachmentUsecase model.IAttachmentUsecase
}

func NewAttachmentHandler(e *echo.Echo, attachmentUsecase model.IAttachmentUsecase) {
	handler := &AttachmentHandler{attachmentUsecase: attachmentUsecase}

	routeUrl := e.Group("v1/attachment")
	routeUrl.POST("/upload", handler.Upload, AuthMiddleware)
	routeUrl.GET("/:ticket_id", handler.FindAllByTicketID, AuthMiddleware)
}

func (h *AttachmentHandler) FindAllByTicketID(ctx echo.Context) error {
	ticketID, err := strconv.ParseInt(ctx.Param("ticket_id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ticket_id")
	}

	attachments, err := h.attachmentUsecase.FindAllByTicketID(ctx.Request().Context(), ticketID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch attachments")
	}

	return ctx.JSON(http.StatusOK, Response{
		Status: http.StatusOK,
		Data:   attachments,
	})
}

func (h *AttachmentHandler) Upload(ctx echo.Context) error {
	ticketIDStr := ctx.FormValue("ticket_id")
	ticketID, err := strconv.ParseInt(ticketIDStr, 10, 64)
	if err != nil || ticketID == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ticket_id")
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "file is required")
	}

	if _, err := os.Stat("./uploads/tickets"); os.IsNotExist(err) {
		os.MkdirAll("./uploads/tickets", os.ModePerm)
	}

	fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(file.Filename))
	savePath := filepath.Join("./uploads/tickets", fileName)

	src, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to open uploaded file")
	}
	defer src.Close()

	dst, err := os.Create(savePath)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create destination file")
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to save file")
	}

	input := model.CreateAttachmentInput{
		TicketID: ticketID,
		FilePath: savePath,
	}

	err = h.attachmentUsecase.Create(ctx.Request().Context(), input)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create attachment")
	}

	return ctx.JSON(http.StatusCreated, Response{
		Status:  http.StatusCreated,
		Message: "Attachment created successfully",
	})
}
