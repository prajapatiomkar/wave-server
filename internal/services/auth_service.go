package services

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/prajapatiomkar/wave-server/internal/models"
	"github.com/prajapatiomkar/wave-server/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *repositories.UserRepository
}

func NewAuthService(userRepo *repositories.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	User  models.UserResponse `json:"user"`
	Token string              `json:"token"`
}

func (s *AuthService) Register(req *RegisterRequest) (*AuthResponse, error) {
	if _, err := s.userRepo.FindByEmail(req.Email); err == nil {
		return nil, errors.New("email already registered")
	}

	if _, err := s.userRepo.FindByUsername(req.Username); err == nil {
		return nil, errors.New("username already taken")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		FullName: req.FullName,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("failed to create user")
	}

	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &AuthResponse{
		User:  s.toUserResponse(user),
		Token: token,
	}, nil
}

func (s *AuthService) Login(req *LoginRequest) (*AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &AuthResponse{
		User:  s.toUserResponse(user),
		Token: token,
	}, nil
}

func (s *AuthService) GetUserByID(userID uint) (*models.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	userResp := s.toUserResponse(user)
	return &userResp, nil
}

func (s *AuthService) generateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func (s *AuthService) toUserResponse(user *models.User) models.UserResponse {
	return models.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		FullName: user.FullName,
		Avatar:   user.Avatar,
		IsOnline: user.IsOnline,
		LastSeen: user.LastSeen,
	}
}
