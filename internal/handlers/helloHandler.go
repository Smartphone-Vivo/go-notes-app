package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type HelloHandler struct {
	semaphore chan struct{}
}

func NewHelloHandler() *HelloHandler {
	return &HelloHandler{
		semaphore: make(chan struct{}, 3),
	}
}

func (h *HelloHandler) HelloEndpoint(c echo.Context) error {
	select {
	case h.semaphore <- struct{}{}:

		defer func() { <-h.semaphore }()

		time.Sleep(1 * time.Second)
		return c.JSON(http.StatusOK, "hello")

	default:
		return c.JSON(http.StatusBadRequest, "goodbye")
	}

}
