package response

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

var (
	ErrInvalidInput  = errors.New("invalid input")
	ErrNotFound      = errors.New("resource not found")
	ErrInsufficient  = errors.New("insufficient balance")
	ErrAlreadyExists = errors.New("resource already exists")
)

func InternalServerError(c echo.Context, err error) error {
	log.Error().
		Err(err).
		Str("method", c.Request().Method).
		Str("path", c.Path()).
		Msg("internal server error")

	return c.JSON(http.StatusInternalServerError, echo.Map{
		"error": "the server encountered a problem and could not process your request",
	})
}

func BadRequest(c echo.Context, err error) error {
	log.Warn().
		Err(err).
		Str("method", c.Request().Method).
		Str("path", c.Path()).
		Msg("bad request")

	return c.JSON(http.StatusBadRequest, echo.Map{
		"error": err.Error(),
	})
}

func NotFound(c echo.Context, err error) error {
	log.Warn().
		Err(err).
		Str("method", c.Request().Method).
		Str("path", c.Path()).
		Msg("not found")

	return c.JSON(http.StatusNotFound, echo.Map{
		"error": "not found",
	})
}

func UnprocessableEntity(c echo.Context, err error) error {
	log.Warn().
		Err(err).
		Str("method", c.Request().Method).
		Str("path", c.Path()).
		Msg("unprocessable entity")

	return c.JSON(http.StatusUnprocessableEntity, echo.Map{
		"error": err.Error(),
	})
}

func Unauthorized(c echo.Context, err error) error {
	log.Warn().
		Err(err).
		Str("method", c.Request().Method).
		Str("path", c.Path()).
		Msg("unauthorized access attempt")

	return c.JSON(http.StatusUnauthorized, echo.Map{
		"error": "invalid or missing authentication token",
	})
}
