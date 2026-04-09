package service

import (
	"context"
	"strings"
	"time"

	"errors"

	"github.com/adishjain1107/tradex/pkg/auth/helper"
	"github.com/adishjain1107/tradex/pkg/auth/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrInvalidPayload     = errors.New("invalid registration payload")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrEmailUppercase     = errors.New("email must be lowercase")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type AuthService struct {
	users *mongo.Collection
}

func NewAuthService(db *mongo.Database) *AuthService {
	return &AuthService{
		users: db.Collection("users"),
	}
}
func (s *AuthService) RegisterService(ctx context.Context, req models.RegReq) (*models.User, error) {
	if s.users == nil {
		return nil, errors.New("Service not initialized")
	}

	email := strings.TrimSpace(req.Email)

	if email == "" || req.Password == "" {
		return nil, ErrInvalidPayload
	}

	if email != strings.ToLower(email) {
		return nil, ErrEmailUppercase
	}

	role := strings.TrimSpace(req.Role)

	if role == "" {
		role = "user"
	}

	var existingUser models.User

	err := s.users.FindOne(ctx, bson.M{"email": email}).Decode(&existingUser)

	if err == nil {
		return nil, ErrUserAlreadyExists
	}
	if !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}

	hash, err := helper.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()

	user := &models.User{
		Email:        email,
		PasswordHash: hash,
		Role:         role,
		Wallet: models.Wallet{
			USD:         0,
			TotalTrades: 0,
		},
		Created: now,
		Updated: now,
	}

	_, err = s.users.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthService) LoginService(ctx context.Context, req models.LoginReq) (*models.AuthResp, error) {
	if s.users == nil {
		return nil, errors.New("Service not initialized")
	}

	email := strings.TrimSpace(req.Email)
	if email == "" || req.Password == "" {
		return nil, ErrInvalidPayload
	}

	if email != strings.ToLower(email) {
		return nil, ErrEmailUppercase
	}

	var existingUser models.User
	err := s.users.FindOne(ctx, bson.M{"email": email}).Decode(&existingUser)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if err := helper.VerifyPassword(existingUser.PasswordHash, req.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	accessToken, err := helper.GenerateAccessToken(existingUser.Email, existingUser.Role)
	if err != nil {
		return nil, err
	}

	refreshToken, err := helper.GenerateRefreshToken(existingUser.Email, existingUser.Role)
	if err != nil {
		return nil, err
	}

	return &models.AuthResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Email:        existingUser.Email,
		Role:         existingUser.Role,
	}, nil
}
