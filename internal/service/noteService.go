package service

import (
	"context"
	"github.com/google/uuid"
	"test-task/internal/domain"
	"test-task/internal/kafka"
	"test-task/internal/repository"
	"time"
)

type NoteService interface {
	CreateNote(ctx context.Context, userID, noteTitle, noteContent string) (domain.Note, error)
	GetAllNotes(ctx context.Context) ([]domain.Note, error)
	DeleteNote(ctx context.Context, id string) error
	GetNoteById(ctx context.Context, id string) (domain.Note, error)
	UpdateNote(ctx context.Context, id string, updateReq domain.UpdateNoteRequest) (domain.Note, error)
}

type noteService struct {
	repo     repository.NotesRepository
	producer *kafka.Producer
}

func NewNoteService(r repository.NotesRepository, p *kafka.Producer) NoteService {
	return &noteService{
		repo:     r,
		producer: p,
	}
}

func (s *noteService) CreateNote(ctx context.Context, userID, noteTitle, noteContent string) (domain.Note, error) {
	note := domain.Note{
		ID:        uuid.NewString(),
		UserID:    userID,
		Title:     noteTitle,
		Content:   noteContent,
		CreatedAt: time.Now(),
	}

	if err := s.producer.PublishNote(ctx, note); err != nil {
		return domain.Note{}, err
	}

	return note, nil
}

func (s *noteService) GetAllNotes(ctx context.Context) ([]domain.Note, error) {
	return s.repo.GetAllNotes(ctx)
}

func (s *noteService) DeleteNote(ctx context.Context, id string) error {
	return s.repo.DeleteNote(ctx, id)
}

func (s *noteService) GetNoteById(ctx context.Context, id string) (domain.Note, error) {
	return s.repo.GetNoteById(ctx, id)
}

func (s *noteService) UpdateNote(ctx context.Context, id string, updateReq domain.UpdateNoteRequest) (domain.Note, error) {
	updatedNote, err := s.repo.GetNoteById(ctx, id)
	if err != nil {
		return domain.Note{}, err
	}

	if updateReq.Title != "" {
		updatedNote.Title = updateReq.Title
	}
	if updateReq.Content != "" {
		updatedNote.Content = updateReq.Content
	}

	if err := s.repo.UpdateNote(ctx, updatedNote); err != nil {
		return domain.Note{}, err
	}

	return updatedNote, nil
}
