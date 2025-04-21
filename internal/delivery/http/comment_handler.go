package http

import (
	"helpdesk-ticketing-system/internal/model"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type CommentHandler struct {
	commentUsecase model.ICommentUsecase
}

func NewCommentHandler(e *echo.Echo, commentUsecase model.ICommentUsecase) {
	handler := &CommentHandler{commentUsecase: commentUsecase}

	routeUrl := e.Group("v1/comment")
	routeUrl.GET("", handler.FindAll, AuthMiddleware)
	routeUrl.GET("/:id", handler.FindById, AuthMiddleware)
	routeUrl.POST("/create", handler.Create, AuthMiddleware)
	routeUrl.PUT("/update/:id", handler.Update, AuthMiddleware)
	routeUrl.DELETE("/delete/:id", handler.Delete, AuthMiddleware)
}

func (c *CommentHandler) FindAll(ctx echo.Context) error {
	comments, err := c.commentUsecase.FindAll(ctx.Request().Context(), model.Comment{})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, Response{
		Status: http.StatusOK,
		Data:   comments,
	})
}

func (c *CommentHandler) FindById(ctx echo.Context) error {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid comment ID format")
	}

	comment, err := c.commentUsecase.FindById(ctx.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Comment not found")
	}

	return ctx.JSON(http.StatusOK, Response{
		Status: http.StatusOK,
		Data:   comment,
	})
}

func (c *CommentHandler) Create(ctx echo.Context) error {
	var body model.CreateCommentInput
	if err := ctx.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	comment, err := c.commentUsecase.Create(ctx.Request().Context(), body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create comment")
	}

	return ctx.JSON(http.StatusCreated, Response{
		Status:  http.StatusCreated,
		Message: "Comment created successfully",
		Data:    comment,
	})
}

func (c *CommentHandler) Update(ctx echo.Context) error {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid comment ID format")
	}

	var body model.UpdateCommentInput
	if err := ctx.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	comment, err := c.commentUsecase.Update(ctx.Request().Context(), id, body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update comment")
	}

	return ctx.JSON(http.StatusOK, Response{
		Status:  http.StatusOK,
		Message: "Comment updated successfully",
		Data:    comment,
	})
}

func (c *CommentHandler) Delete(ctx echo.Context) error {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid comment ID format")
	}

	err = c.commentUsecase.Delete(ctx.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete comment")
	}

	return ctx.JSON(http.StatusOK, Response{
		Status:  http.StatusOK,
		Message: "Comment deleted successfully",
	})
}
