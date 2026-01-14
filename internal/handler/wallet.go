package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/bryanaleron193/wallet-service/internal/service"
	"github.com/bryanaleron193/wallet-service/pkg/utils"
	"github.com/labstack/echo/v4"
)

type WalletHandler struct {
	BaseHandler
	walletService service.WalletService
}

func NewWalletHandler(ws service.WalletService) *WalletHandler {
	return &WalletHandler{walletService: ws}
}

type BalanceResponse struct {
	UserID  string `json:"user_id"`
	Balance string `json:"balance"`
}

func (h *WalletHandler) GetBalance(c echo.Context) error {
	userID := c.Param("user_id")
	if userID == "" {
		return h.badRequestResponse(c, fmt.Errorf("user_id is required"))
	}

	ctx := c.Request().Context()
	wallet, err := h.walletService.GetByUserID(ctx, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return h.notFoundResponse(c, err)
		}

		return h.internalServerError(c, err)
	}

	formattedBalance := utils.FormatRupiah(wallet.Balance)

	return c.JSON(http.StatusOK, BalanceResponse{
		UserID:  wallet.UserID,
		Balance: formattedBalance,
	})
}
