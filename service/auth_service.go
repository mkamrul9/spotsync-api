package service

import (
	"errors"

	"github.com/mkamrul9/spotsync-api/dto"
	"github.com/mkamrul9/spotsync-api/models"
	"github.com/mkamrul9/spotsync-api/repository"
	"github.com/mkamrul9/spotsync-api/utils"
)

type AuthService interface {
	Register(req dto.RegisterRequest) (*dto.RegisterResponse, error)
	Login(req dto.LoginRequest) (*dto.AuthResponse, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo}
}

func (s *authService) Register(req dto.RegisterRequest) (*dto.RegisterResponse, error) {
	// 1. Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// 2. Set default role if empty
	role := req.Role
	if role == "" {
		role = "driver"
	}

	// 3. Create User Model
	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     role,
	}

	// 4. Save via Repository
	if err := s.userRepo.CreateUser(&user); err != nil {
		return nil, errors.New("email already exists or database error")
	}

	return &dto.RegisterResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *authService) Login(req dto.LoginRequest) (*dto.AuthResponse, error) {
	// 1. Fetch user by email
	user, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// 2. Check password
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	// 3. Generate JWT
	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &dto.AuthResponse{
		Token: token,
		User: dto.UserSummary{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
	}, nil
}
