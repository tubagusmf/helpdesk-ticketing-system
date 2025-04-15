package http

import (
	"context"
	"log"
	"net/http"
	"strings"

	"helpdesk-ticketing-system/internal/helper"
	"helpdesk-ticketing-system/internal/model"

	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get(echo.HeaderAuthorization)
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing token")
		}

		splitAuth := strings.Split(authHeader, " ")
		if len(splitAuth) != 2 || splitAuth[0] != "Bearer" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token format")
		}

		accessToken := splitAuth[1]

		var claim model.CustomClaims
		err := helper.DecodeToken(accessToken, &claim)
		if err != nil {
			log.Println("Token decoding failed:", err)
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid or expired token")
		}

		ctx := context.WithValue(c.Request().Context(), model.BearerAuthKey, claim)
		req := c.Request().WithContext(ctx)
		c.SetRequest(req)

		return next(c)
	}
}
