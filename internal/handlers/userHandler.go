package handlers

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"test-task/internal/appErrors"
	"test-task/internal/domain"
	"test-task/internal/service"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) handleError(c echo.Context, err error) error {
	switch {
	case errors.As(err, &appErrors.ValidationError{}):
		var valErr appErrors.ValidationError
		errors.As(err, &valErr)

		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": valErr.Error()})

	case errors.As(err, &appErrors.NotFoundError{}):
		var notFound appErrors.NotFoundError
		errors.As(err, &notFound)

		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})

	case errors.As(err, &appErrors.DatabaseError{}):
		var dbErr appErrors.DatabaseError
		errors.As(err, &dbErr)

		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal server error",
		})

	default:
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal server error",
		})
	}
}

func (h *UserHandler) Register(c echo.Context) error {
	var req domain.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request",
		})
	}

	resp, err := h.userService.Register(c.Request().Context(), req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(http.StatusCreated, resp)
}

func (h *UserHandler) Login(c echo.Context) error {
	var req domain.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request",
		})
	}

	resp, err := h.userService.Login(c.Request().Context(), req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, resp)
}
