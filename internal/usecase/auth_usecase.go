package usecase

import (
	"context"
	"errors"
	"fmt"
	"os"
	"project-bcc/dto"
	"project-bcc/internal/entity"
	"project-bcc/pkg/bcrypt"
	"project-bcc/pkg/email"
	"project-bcc/pkg/jwt"

	"gorm.io/gorm"
)

var (
	ErrConflictData       = errors.New("Email sudah digunakan")
	ErrInvalidCredentials = errors.New("Email atau password salah")
	ErrUnauthorized       = errors.New("Token tidak valid atau sudah kedaluwarsa")
	ErrUnverifiedAccount  = errors.New("Akun belum diverifikasi, silakan cek email anda")
	ErrInternalServer     = errors.New("Terjadi masalah internal pada server")
	ErrAlreadyVerified    = errors.New("Email sudah diverifikasi")
)

type AuthUsecase struct {
	authRepo AuthRepository
}

func NewAuthUsecase(repo AuthRepository) *AuthUsecase {
	return &AuthUsecase{authRepo: repo}
}

func (au *AuthUsecase) Register(ctx context.Context, req dto.RegisterRequest) error {
	existEmail, _ := au.authRepo.FindByEmail(ctx, req.Email)
	if existEmail != nil {
		return ErrConflictData
	}

	hashedPassword, err := bcrypt.HashPassword(req.Password)
	if err != nil {
		fmt.Println("gagal membuat password")
		return ErrInternalServer
	}

	userRole := entity.RoleUser

	user := &entity.User{
		FullName:   req.Fullname,
		Email:      req.Email,
		Password:   hashedPassword,
		Role:       userRole,
		IsVerified: false,
	}

	err = au.authRepo.Save(ctx, user)
	if err != nil {
		fmt.Println("Gagal menyimpan user" + err.Error())
		return ErrInternalServer
	}

	frontendURL := os.Getenv("FE_URL")
	token, _ := jwt.GenerateEmailVerificationToken(user.ID.String())
	verificationLink := fmt.Sprintf("%s/verify", frontendURL)
	err = email.SendVerificationEmail(user.Email, verificationLink, token)
	if err != nil {
		fmt.Println("Gagal mengirim email verifikasi" + err.Error())
		return ErrInternalServer
	}

	return nil
}

func (au *AuthUsecase) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {

	user, err := au.authRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("User tidak ditemukan" + err.Error())
			return nil, ErrInvalidCredentials
		}
		return nil, ErrInternalServer
	}

	if !user.IsVerified {
		return nil, ErrUnverifiedAccount
	}

	if !bcrypt.CheckPassword(user.Password, req.Password) {
		return nil, ErrInvalidCredentials
	}

	accessToken, err := jwt.GenerateAccessToken(user.ID.String(), string(user.Role), user.UpdatedAt.Unix())
	if err != nil {
		fmt.Println("Gagal generate access token" + err.Error())
		return nil, ErrInternalServer
	}

	refreshToken, err := jwt.GenerateRefreshToken(user.ID.String(), string(user.Role), user.UpdatedAt.Unix(), req.RememberMe)
	if err != nil {
		fmt.Println("Gagal generate refresh token" + err.Error())
		return nil, ErrInternalServer
	}

	err = au.authRepo.UpdateRefreshToken(ctx, user.ID.String(), refreshToken)
	if err != nil {
		fmt.Println("Gagal menyimpan refresh token" + err.Error())
		return nil, ErrInternalServer
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.UserData{
			ID:       user.ID.String(),
			Fullname: user.FullName,
			Email:    user.Email,
			Role:     string(user.Role),
		},
	}, nil
}

func (au *AuthUsecase) VerifyEmail(ctx context.Context, req dto.VerifyEmailRequest) error {

	claims, err := jwt.ValidateToken(req.Token, os.Getenv("EMAIL_VERIFY_SECRET"))
	if err != nil {
		fmt.Println("Token verifikasi tidak valid atau sudah kedaluwarsa" + err.Error())
		return ErrUnauthorized
	}

	if claims.Type != "email_verification" {
		fmt.Println("Tipe token tidak valid")
		return ErrUnauthorized
	}

	user, err := au.authRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("User tidak ditemukan" + err.Error())
			return ErrInvalidCredentials
		}
		return ErrInternalServer
	}

	if user.IsVerified {
		return ErrAlreadyVerified
	}

	user.IsVerified = true
	err = au.authRepo.Update(ctx, user)
	if err != nil {
		fmt.Println("Gagal memperbarui status verifikasi email" + err.Error())
		return ErrInternalServer
	}
	return nil
}

func (au *AuthUsecase) ResendVerification(ctx context.Context, req dto.ResendVerificationRequest) error {

	user, err := au.authRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("User tidak ditemukan" + err.Error())
			return ErrInvalidCredentials
		}
		return ErrInternalServer
	}

	if user.IsVerified {
		return ErrAlreadyVerified
	}

	token, err := jwt.GenerateEmailVerificationToken(user.ID.String())
	if err != nil {
		fmt.Println("Gagal generate token" + err.Error())
		return ErrInternalServer
	}

	link := req.CallbackUrl
	if link == "" {
		link = os.Getenv("FE_URL") + "/verify"
	}

	err = email.SendVerificationEmail(user.Email, link, token)
	if err != nil {
		fmt.Println("Gagal mengirim email verifikasi" + err.Error())
		return ErrInternalServer
	}
	return nil
}

func (au *AuthUsecase) ForgotPassword(ctx context.Context, req dto.ForgotPasswordRequest) error {

	user, err := au.authRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("User tidak ditemukan" + err.Error())
			return ErrInvalidCredentials
		}
		return ErrInternalServer
	}

	token, err := jwt.GenerateResetPasswordToken(user.ID.String(), user.UpdatedAt.Unix())
	if err != nil {
		fmt.Println("Gagal generate token" + err.Error())
		return ErrInternalServer
	}

	link := os.Getenv("FE_URL") + "/reset-password"

	err = email.SendResetPasswordEmail(user.Email, link, token)
	if err != nil {
		fmt.Println("Gagal mengirim email reset password" + err.Error())
		return ErrInternalServer
	}
	return nil
}

func (au *AuthUsecase) ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) error {

	claims, err := jwt.ValidateToken(req.Token, os.Getenv("RESET_PASSWORD_SECRET"))
	if err != nil {
		fmt.Println("Token tidak valid atau sudah kedaluwarsa" + err.Error())
		return ErrUnauthorized
	}

	if claims.Type != "reset_password" {
		fmt.Println("Tipe token tidak valid")
		return ErrUnauthorized
	}

	user, err := au.authRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("User tidak ditemukan" + err.Error())
			return ErrInvalidCredentials
		}
		return ErrInternalServer
	}

	if claims.UpdatedAt < user.UpdatedAt.Unix() {
		fmt.Println(ErrUnauthorized)
		return ErrUnauthorized
	}

	hashedPassword, err := bcrypt.HashPassword(req.Password)
	if err != nil {
		fmt.Println("Gagal hash password" + err.Error())
		return ErrInternalServer
	}

	user.Password = hashedPassword
	user.RefreshToken = ""
	err = au.authRepo.Update(ctx, user)
	if err != nil {
		fmt.Println("Gagal memperbarui password" + err.Error())
		return ErrInternalServer
	}
	return nil
}

func (au *AuthUsecase) Refresh(ctx context.Context, refreshToken string) (*dto.RefreshResponse, error) {

	claims, err := jwt.ValidateToken(refreshToken, os.Getenv("REFRESH_TOKEN_SECRET"))
	if err != nil {
		fmt.Println("Refresh token tidak valid atau kedaluwarsa" + err.Error())
		return nil, ErrUnauthorized
	}

	if claims.Type != "refresh" {
		fmt.Println("Tipe token tidak valid")
		return nil, ErrUnauthorized
	}

	user, err := au.authRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("User tidak ditemukan" + err.Error())
			return nil, ErrInvalidCredentials
		}
		return nil, ErrInternalServer
	}

	if user.RefreshToken != refreshToken {
		fmt.Println("Refresh token tidak sesuai")
		return nil, ErrUnauthorized
	}

	accessToken, err := jwt.GenerateAccessToken(user.ID.String(), string(user.Role), user.UpdatedAt.Unix())
	if err != nil {
		fmt.Println("Gagal generate access token" + err.Error())
		return nil, ErrInternalServer
	}

	return &dto.RefreshResponse{
		AccessToken: accessToken,
	}, nil
}

func (au *AuthUsecase) Logout(ctx context.Context, refreshToken string) error {
	claims, err := jwt.ValidateToken(refreshToken, os.Getenv("REFRESH_TOKEN_SECRET"))
	if err != nil {
		fmt.Println("Refresh token tidak valid atau kedaluwarsa" + err.Error())
		return ErrUnauthorized
	}

	user, err := au.authRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("User tidak ditemukan" + err.Error())
			return ErrInvalidCredentials
		}
		return ErrInternalServer
	}

	user.RefreshToken = ""
	err = au.authRepo.Update(ctx, user)
	if err != nil {
		fmt.Println("Gagal mengupdate refresh token" + err.Error())
		return ErrInternalServer
	}
	return nil
}
