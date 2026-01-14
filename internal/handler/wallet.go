package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/bryanaleron193/wallet-service/internal/service"
	"github.com/bryanaleron193/wallet-service/pkg/response"
	"github.com/bryanaleron193/wallet-service/pkg/util"
	"github.com/labstack/echo/v4"
)

type WalletHandler struct {
	walletService service.WalletService
}

func NewWalletHandler(ws service.WalletService) *WalletHandler {
	return &WalletHandler{walletService: ws}
}

type BalanceResponse struct {
	UserID  string `json:"user_id"`
	Balance string `json:"balance"`
}

// GetBalance godoc
//
//	@Summary		Get Wallet Balance
//	@Description	Get user's balance
//	@Tags			Wallet
//	@Security		BearerAuth
//	@Success		200	{object}	BalanceResponse
//	@Router			/api/v1/wallets/balance [get]
func (h *WalletHandler) GetBalance(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return response.Unauthorized(c, nil)
	}

	ctx := c.Request().Context()
	wallet, err := h.walletService.GetByUserID(ctx, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return response.NotFound(c, err)
		}

		return response.InternalServerError(c, err)
	}

	formattedBalance := util.FormatRupiah(wallet.Balance)

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

// Withdraw godoc
//
//	@Summary		Withdraw Money
//	@Description	Withdraw money from user's account
//	@Tags			Wallet
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		WithdrawRequest	true	"Withdrawal Amount"
//	@Success		200		{object}	WithdrawResponse
//	@Router			/api/v1/wallets/withdraw [post]
func (h *WalletHandler) Withdraw(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return response.Unauthorized(c, nil)
	}

	var req WithdrawRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, fmt.Errorf("invalid request payload"))
	}

	ctx := c.Request().Context()

	updatedWallet, transactionID, err := h.walletService.Withdraw(ctx, userID, req.Amount, req.Description)

	if err != nil {
		if errors.Is(err, response.ErrInsufficient) {
			return response.UnprocessableEntity(c, fmt.Errorf("insufficient balance"))
		}
		if errors.Is(err, response.ErrNotFound) {
			return response.NotFound(c, err)
		}
		return response.InternalServerError(c, err)
	}

	return c.JSON(http.StatusOK, WithdrawResponse{
		TransactionID:  transactionID,
		UserID:         userID,
		Amount:         util.FormatRupiah(req.Amount),
		CurrentBalance: util.FormatRupiah(updatedWallet.Balance),
		Message:        "withdrawal successful",
	})
}
