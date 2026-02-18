package repository

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"test-task/internal/appErrors"
	"test-task/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByUsername(ctx context.Context, username string) (*domain.User, error)
	FindById(ctx context.Context, id string) (*domain.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return appErrors.DatabaseError{Err: err, Op: "hash password"}
	}

	user.Password = string(hashedPassword)

	err = r.db.WithContext(ctx).Create(user).Error
	if err != nil {
		return appErrors.DatabaseError{Err: err, Op: "create user"}
	}

	return nil
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user domain.User

	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, appErrors.NotFoundError{Entity: "User", ID: username}
	}
	if err != nil {
		return nil, appErrors.DatabaseError{Err: err, Op: "FindByUsername"}
	}
	return &user, nil
}

func (r *userRepository) FindById(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, appErrors.NotFoundError{Entity: "User", ID: id}
	}
	if err != nil {
		return nil, appErrors.DatabaseError{Err: err, Op: "FindById"}
	}

	return &user, nil
}
