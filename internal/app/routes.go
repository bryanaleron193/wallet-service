package app

import (
	"github.com/bryanaleron193/wallet-service/internal/middleware"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func (c *Container) RegisterRoutes(e *echo.Echo) {
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.POST("/login", c.UserHandler.Login)

	api := e.Group("/api/v1")

	api.Use(middleware.AuthMiddleware(c.Config.JWT.Secret))

	wallets := api.Group("/wallets")
	{
		wallets.GET("/balance", c.WalletHandler.GetBalance)
		wallets.POST("/withdraw", c.WalletHandler.Withdraw)
	}
}
