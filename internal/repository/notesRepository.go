package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"test-task/internal/appErrors"
	"test-task/internal/domain"
)

type NotesRepository interface {
	CreateNote(ctx context.Context, note domain.Note) error
	GetAllNotes(ctx context.Context) ([]domain.Note, error)
	GetNoteById(ctx context.Context, id string) (domain.Note, error)
	UpdateNote(ctx context.Context, note domain.Note) error
	DeleteNote(ctx context.Context, id string) error
}

type noteRepository struct {
	db *gorm.DB
}

func NewNotesRepository(db *gorm.DB) NotesRepository {
	return &noteRepository{db: db}
}

func (r *noteRepository) CreateNote(ctx context.Context, note domain.Note) error {

	err := r.db.WithContext(ctx).Create(&note).Error

	if err != nil {
		return appErrors.DatabaseError{Err: err, Op: "create"}
	}

	return nil
}

func (r *noteRepository) GetAllNotes(ctx context.Context) ([]domain.Note, error) {
	var allNotes []domain.Note
	err := r.db.WithContext(ctx).Find(&allNotes).Error

	if err != nil {
		return nil, appErrors.DatabaseError{Err: err, Op: "error gtt all"}
	}

	return allNotes, err
}

func (r *noteRepository) GetNoteById(ctx context.Context, id string) (domain.Note, error) {
	var note domain.Note
	err := r.db.WithContext(ctx).First(&note, "id = ?", id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.Note{}, appErrors.NotFoundError{
			Entity: "Note",
			ID:     "id",
		}
	}

	if err != nil {
		return domain.Note{}, appErrors.DatabaseError{Err: err, Op: "get by id"}
	}

	return note, err
}

func (r *noteRepository) DeleteNote(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&domain.Note{}, "id = ?", id)

	if result.Error != nil {
		return appErrors.DatabaseError{Err: result.Error, Op: "delete"}
	}

	if result.RowsAffected == 0 {
		return appErrors.NotFoundError{
			Entity: "Note",
			ID:     id,
		}
	}

	return nil
}

func (r *noteRepository) UpdateNote(ctx context.Context, note domain.Note) error {

	err := r.db.WithContext(ctx).Save(&note).Error

	if err != nil {
		return appErrors.DatabaseError{Err: err}
	}

	return nil
}
