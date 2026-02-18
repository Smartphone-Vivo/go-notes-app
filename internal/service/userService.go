package service

import (
	"context"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"test-task/internal/appErrors"
	"test-task/internal/auth"
	"test-task/internal/domain"
	"test-task/internal/repository"
)

type UserService interface {
	Register(ctx context.Context, req domain.RegisterRequest) (*domain.AuthResponse, error)
	Login(ctx context.Context, req domain.LoginRequest) (*domain.AuthResponse, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}
func (s *userService) Register(ctx context.Context, req domain.RegisterRequest) (*domain.AuthResponse, error) {
	existing, _ := s.userRepo.FindByUsername(ctx, req.Username)
	if existing != nil {
		return nil, appErrors.ValidationError{
			Field:   "username",
			Message: "username already exists",
		}
	}

	user := &domain.User{
		ID:       uuid.NewString(),
		Username: req.Username,
		Password: req.Password,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		return nil, appErrors.DatabaseError{Err: err, Op: "generate token"}
	}

	return &domain.AuthResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (s *userService) Login(ctx context.Context, req domain.LoginRequest) (*domain.AuthResponse, error) {
	user, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, appErrors.ValidationError{
			Field:   "username",
			Message: "invalid credentials",
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, appErrors.ValidationError{
			Field:   "password",
			Message: "invalid credentials",
		}
	}

	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		return nil, appErrors.DatabaseError{Err: err, Op: "generate token"}
	}
	return &domain.AuthResponse{
		Token: token,
		User:  *user,
	}, nil
}
