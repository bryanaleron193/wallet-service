package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/bryanaleron193/wallet-service/internal/service"
	"github.com/bryanaleron193/wallet-service/pkg/apperror"
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

type WithdrawRequest struct {
	Amount      float64 `json:"amount" validate:"required,gt=0"`
	Description string  `json:"description"`
}

type WithdrawResponse struct {
	TransactionID  string `json:"transaction_id"`
	UserID         string `json:"user_id"`
	Amount         string `json:"amount"`
	CurrentBalance string `json:"current_balance"`
	Message        string `json:"message"`
}

func (h *WalletHandler) Withdraw(c echo.Context) error {
	userID := c.Param("user_id")
	if userID == "" {
		return h.badRequestResponse(c, fmt.Errorf("user_id is required"))
	}

	var req WithdrawRequest
	if err := c.Bind(&req); err != nil {
		return h.badRequestResponse(c, fmt.Errorf("invalid request payload"))
	}

	ctx := c.Request().Context()

	updatedWallet, transactionID, err := h.walletService.Withdraw(ctx, userID, req.Amount, req.Description)

	if err != nil {
		if errors.Is(err, apperror.ErrInsufficient) {
			return h.unprocessableEntityResponse(c, fmt.Errorf("insufficient balance"))
		}
		if errors.Is(err, apperror.ErrNotFound) {
			return h.notFoundResponse(c, err)
		}
		return h.internalServerError(c, err)
	}

	return c.JSON(http.StatusOK, WithdrawResponse{
		TransactionID:  transactionID,
		UserID:         userID,
		Amount:         utils.FormatRupiah(req.Amount),
		CurrentBalance: utils.FormatRupiah(updatedWallet.Balance),
		Message:        "withdrawal successful",
	})
}
