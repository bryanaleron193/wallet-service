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

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

// Login godoc
//
//	@Summary		User Login
//	@Description	Authenticate user and return JWT token
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			request	body		LoginRequest	true	"Login Request"
//	@Success		200		{object}	LoginResponse
//	@Router			/login [post]
func (h *UserHandler) Login(c echo.Context) error {
	var req LoginRequest

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

	return c.JSON(http.StatusOK, LoginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
	})
}
