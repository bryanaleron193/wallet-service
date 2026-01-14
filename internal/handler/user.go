package handler

import (
	"net/http"

	"github.com/bryanaleron193/wallet-service/internal/service"
	"github.com/bryanaleron193/wallet-service/pkg/response"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) Login(c echo.Context) error {
	var req struct {
		Username string `json:"username" validate:"required"`
	}

	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, err)
	}

	if req.Username == "" {
		return response.UnprocessableEntity(c, nil)
	}

	token, err := h.userService.Login(c.Request().Context(), req.Username)
	if err != nil {
		return response.Unauthorized(c, err)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"access_token": "Bearer" + token,
	})
}
