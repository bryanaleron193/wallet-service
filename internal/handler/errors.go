package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type BaseHandler struct{}

func (h *BaseHandler) internalServerError(c echo.Context, err error) error {
	log.Error().
		Err(err).
		Str("method", c.Request().Method).
		Str("path", c.Path()).
		Msg("internal server error")

	return c.JSON(http.StatusInternalServerError, echo.Map{
		"error": "the server encountered a problem and could not process your request",
	})
}

func (h *BaseHandler) badRequestResponse(c echo.Context, err error) error {
	log.Warn().
		Err(err).
		Str("method", c.Request().Method).
		Str("path", c.Path()).
		Msg("bad request")

	return c.JSON(http.StatusBadRequest, echo.Map{
		"error": err.Error(),
	})
}

func (h *BaseHandler) notFoundResponse(c echo.Context, err error) error {
	log.Warn().
		Err(err).
		Str("method", c.Request().Method).
		Str("path", c.Path()).
		Msg("not found")

	return c.JSON(http.StatusNotFound, echo.Map{
		"error": "not found",
	})
}
