package usecase

import (
	"context"
	"errors"
	"project-bcc/dto"
	"project-bcc/internal/entity"
	"project-bcc/pkg/bcrypt"
	"project-bcc/pkg/jwt"

	"github.com/google/uuid"
)

type AuthUsecase struct {
	authRepo AuthRepository
}

func NewAuthUsecase(repo AuthRepository) *AuthUsecase {
	return &AuthUsecase{authRepo: repo}
}

func (au *AuthUsecase) Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error) {
	existEmail, _ := au.authRepo.FindByEmail(ctx, req.Email)
	if existEmail != nil {
		return nil, errors.New("Email sudah digunakan")
	}

	existUsername, _ := au.authRepo.FindByUsername(ctx, req.Username)
	if existUsername != nil {
		return nil, errors.New("Username sudah digunakan")
	}

	hashedPassword, err := bcrypt.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("Gagal Membuat Password")
	}

	user := &entity.User{
		ID:       uuid.New(),
		FullName: req.Fullname,
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     "user",
	}

	err = au.authRepo.Save(ctx, user)
	if err != nil {
		return nil, errors.New("Gagal menyimpan user")
	}

	token, err := jwt.GenerateToken(user.ID.String(), user.Role)
	if err != nil {
		return nil, errors.New("Gagal membuat token")
	}

	return &dto.AuthResponse{
		Token: token,
		UserData: dto.UserData{
			ID:       user.ID.String(),
			Fullname: user.FullName,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
		},
	}, nil
}

func (au *AuthUsecase) Login(ctx context.Context, log dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := au.authRepo.FindByEmail(ctx, log.Email)
	if err != nil || user == nil {
		return nil, errors.New("Email atau password salah")
	}

	if !bcrypt.CheckPassword(user.Password, log.Password) {
		return nil, errors.New("Email atau pasword salah")
	}

	token, err := jwt.GenerateToken(user.ID.String(), user.Role)
	if err != nil {
		return nil, errors.New("Gagal generate token")
	}

	return &dto.AuthResponse{
		Token: token,
		UserData: dto.UserData{
			ID:       user.ID.String(),
			Fullname: user.FullName,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
		},
	}, nil
}
