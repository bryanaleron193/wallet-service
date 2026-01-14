package middleware

import (
	"strings"

	"github.com/bryanaleron193/wallet-service/pkg/response"
	"github.com/bryanaleron193/wallet-service/pkg/util"
	"github.com/labstack/echo/v4"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func AuthMiddleware(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return response.Unauthorized(c, nil)
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			claims, err := util.ValidateToken(tokenString, secret)
			if err != nil {
				return response.Unauthorized(c, err)
			}

			c.Set("user_id", claims.UserID)
			c.Set("username", claims.Username)

			return next(c)
		}
	}
}
