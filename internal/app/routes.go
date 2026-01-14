package app

import "github.com/labstack/echo/v4"

func (c *Container) RegisterRoutes(e *echo.Echo) {
	api := e.Group("/api/v1")

	wallets := api.Group("/wallets")
	{
		wallets.GET("/:user_id/balance", c.WalletHandler.GetBalance)
		wallets.POST("/:user_id/withdraw", c.WalletHandler.Withdraw)
	}
}
