package handlers

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"test-task/internal/appErrors"
	"test-task/internal/domain"
	"test-task/internal/service"
)

type NotesHandler struct {
	service service.NoteService
}

func NewNotesHandlers(s service.NoteService) *NotesHandler {
	return &NotesHandler{service: s}
}

func (h *NotesHandler) handleError(c echo.Context, err error) error {
	switch {

	case errors.As(err, &appErrors.NotFoundError{}):
		var notFound appErrors.NotFoundError
		errors.As(err, &notFound)
		return c.JSON(http.StatusNotFound, map[string]string{"error": notFound.Error()})

	case errors.As(err, &appErrors.ValidationError{}):
		var valErr appErrors.ValidationError
		errors.As(err, &valErr)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": valErr.Error()})

	case errors.As(err, &appErrors.DatabaseError{}):
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})

	default:
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})

	}
}

func (h *NotesHandler) GetNotes(c echo.Context) error {
	ctx := c.Request().Context()

	notes, err := h.service.GetAllNotes(ctx)

	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, notes)
}

func (h *NotesHandler) PostNotes(c echo.Context) error {
	ctx := c.Request().Context()

	var req domain.CreateNoteRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Kal request"})
	}

	note, err := h.service.CreateNote(ctx, req.Title, req.Content)

	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, note)
}

func (h *NotesHandler) DeleteNote(c echo.Context) error {

	ctx := c.Request().Context()

	id := c.Param("id")

	if err := h.service.DeleteNote(ctx, id); err != nil {
		return h.handleError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *NotesHandler) UpdateNote(c echo.Context) error {
	ctx := c.Request().Context()

	id := c.Param("id")

	var req domain.UpdateNoteRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	updatedReq, err := h.service.UpdateNote(ctx, id, req)

	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, updatedReq)
}

func (h *NotesHandler) GetNoteById(c echo.Context) error {
	ctx := c.Request().Context()

	id := c.Param("id")

	note, err := h.service.GetNoteById(ctx, id)

	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, note)

}
