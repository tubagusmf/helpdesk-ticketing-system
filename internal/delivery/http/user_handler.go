package http

import (
	"log"
	"net/http"
	"strconv"

	"helpdesk-ticketing-system/internal/model"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userUsecase model.IUserUsecase
}

func NewUserHandler(e *echo.Echo, userUsecase model.IUserUsecase) {
	handlers := &UserHandler{
		userUsecase: userUsecase,
	}

	routeUser := e.Group("v1/auth")
	routeUser.POST("/login", handlers.Login)
	routeUser.POST("/logout", handlers.Logout, AuthMiddleware)
	routeUser.GET("/user/:id", handlers.FindById, AuthMiddleware)
	routeUser.GET("/users", handlers.FindAll, AuthMiddleware)
	routeUser.POST("/register", handlers.Create)
	routeUser.PUT("/update/:id", handlers.Update, AuthMiddleware)
	routeUser.DELETE("/delete/:id", handlers.Delete, AuthMiddleware)
}

func (handler *UserHandler) Login(c echo.Context) error {
	var body model.LoginInput
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	accessToken, err := handler.userUsecase.Login(c.Request().Context(), body)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid Email or Password")
	}

	return c.JSON(http.StatusOK, Response{
		Status:      http.StatusOK,
		Message:     "Login successful",
		AccessToken: accessToken,
	})
}

func (handler *UserHandler) Logout(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing token")
	}

	err := handler.userUsecase.Logout(c.Request().Context(), token)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to logout")
	}

	return c.JSON(http.StatusOK, Response{
		Status:  http.StatusOK,
		Message: "Logout successful",
	})
}

func (handler *UserHandler) FindById(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID format")
	}

	claim, ok := c.Request().Context().Value(model.BearerAuthKey).(model.CustomClaims)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	if claim.UserID != id {
		return echo.NewHTTPError(http.StatusForbidden, "Access denied")
	}

	log.Printf("Authenticated User ID: %d", claim.UserID)

	user, err := handler.userUsecase.FindById(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	return c.JSON(http.StatusOK, Response{
		Status: http.StatusOK,
		Data:   user,
	})
}

func (handler *UserHandler) FindAll(c echo.Context) error {
	var filter model.User
	filter.Name = c.QueryParam("name")
	filter.Email = c.QueryParam("email")

	users, err := handler.userUsecase.FindAll(c.Request().Context(), filter)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch users")
	}

	return c.JSON(http.StatusOK, Response{
		Status: http.StatusOK,
		Data:   users,
	})
}

func (handler *UserHandler) Create(c echo.Context) error {
	var body model.CreateUserInput
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if body.Name == "" || body.Email == "" || body.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "All fields are required")
	}

	accessToken, err := handler.userUsecase.Create(c.Request().Context(), body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create user")
	}

	return c.JSON(http.StatusCreated, Response{
		Status:      http.StatusCreated,
		Message:     "User registered successfully",
		AccessToken: accessToken,
	})
}

func (handler *UserHandler) Update(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID format")
	}

	claim, ok := c.Request().Context().Value(model.BearerAuthKey).(model.CustomClaims)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	// User hanya boleh mengupdate datanya sendiri
	if claim.UserID != id {
		return echo.NewHTTPError(http.StatusForbidden, "Access denied")
	}

	var body model.UpdateUserInput
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	err = handler.userUsecase.Update(c.Request().Context(), id, body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update user")
	}

	return c.JSON(http.StatusOK, Response{
		Status:  http.StatusOK,
		Message: "User updated successfully",
	})
}

func (handler *UserHandler) Delete(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID format")
	}

	claim, ok := c.Request().Context().Value(model.BearerAuthKey).(model.CustomClaims)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	if claim.UserID != id {
		return echo.NewHTTPError(http.StatusForbidden, "Access denied")
	}

	err = handler.userUsecase.Delete(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete user")
	}

	return c.JSON(http.StatusOK, Response{
		Status:  http.StatusOK,
		Message: "User deleted successfully",
	})
}
