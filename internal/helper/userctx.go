package helper

import (
	"context"
	"errors"
	"helpdesk-ticketing-system/internal/model"
)

func GetUserID(ctx context.Context) (int64, error) {
	val := ctx.Value(model.BearerAuthKey)
	if val == nil {
		return 0, errors.New("user claims not found in context")
	}

	claims, ok := val.(model.CustomClaims)
	if !ok {
		return 0, errors.New("invalid claims type in context")
	}

	return claims.UserID, nil
}
