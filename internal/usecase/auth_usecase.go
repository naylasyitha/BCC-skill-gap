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
)

var (
	ErrInvalidCredentials = errors.New("Email atau password salah")
	ErrUnverifiedAccount  = errors.New("Akun belum diverifikasi, silakan cek email anda")
	ErrInternalServer     = errors.New("Terjadi kesalahan pada server")
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
		return errors.New("Email sudah digunakan")
	}

	hashedPassword, err := bcrypt.HashPassword(req.Password)
	if err != nil {
		return errors.New("Gagal Membuat Password")
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
		return errors.New("Gagal menyimpan user")
	}

	frontendURL := os.Getenv("FE_URL")
	token, _ := jwt.GenerateEmailVerificationToken(user.ID.String())
	fmt.Println("EMAIL VERIFICATION TOKEN: ", token)
	verificationLink := fmt.Sprintf("%s/verify", frontendURL)
	err = email.SendVerificationEmail(user.Email, verificationLink, token)
	if err != nil {
		fmt.Println("EMAIL ERROR: ", err)
		return errors.New("Gagal mengirim email verifikasi" + err.Error())
	}

	return nil
}

func (au *AuthUsecase) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {

	user, err := au.authRepo.FindByEmail(ctx, req.Email)
	if err != nil || user == nil {
		return nil, ErrInvalidCredentials
	}

	if !user.IsVerified {
		return nil, ErrUnverifiedAccount
	}

	if !bcrypt.CheckPassword(user.Password, req.Password) {
		return nil, ErrInvalidCredentials
	}

	accessToken, err := jwt.GenerateAccessToken(user.ID.String(), string(user.Role), user.UpdatedAt.Unix())
	if err != nil {
		return nil, errors.New("Gagal generate access token")
	}

	refreshToken, err := jwt.GenerateRefreshToken(user.ID.String(), string(user.Role), user.UpdatedAt.Unix(), req.RememberMe)
	if err != nil {
		return nil, errors.New("Gagal generate refresh token")
	}

	fmt.Println("ACCESS TOKEN: ", accessToken)
	err = au.authRepo.UpdateRefreshToken(ctx, user.ID.String(), refreshToken)
	if err != nil {
		return nil, errors.New("Gagal menyimpan refresh token")
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
		return errors.New("Token verifikasi tidak valid atau sudah kedaluwarsa")
	}

	if claims.Type != "email_verification" {
		return errors.New("Tipe token tidak valid")
	}

	user, err := au.authRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return errors.New("User tidak ditemukan")
	}

	if user.IsVerified {
		return errors.New("Email sudah diverifikasi")
	}

	user.IsVerified = true
	return au.authRepo.Update(ctx, user)
}

func (au *AuthUsecase) ResendVerification(ctx context.Context, req dto.ResendVerificationRequest) error {

	user, err := au.authRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return errors.New("User tidak ditemukan")
	}

	if user.IsVerified {
		return errors.New("Email sudah diverifikasi")
	}

	token, err := jwt.GenerateEmailVerificationToken(user.ID.String())
	if err != nil {
		return errors.New("Gagal generate token")
	}

	fmt.Println("RESEND EMAIL TOKEN: ", token)

	link := req.CallbackUrl
	if link == "" {
		link = os.Getenv("FE_URL") + "/verify"
	}
	return email.SendVerificationEmail(user.Email, link, token)
}

func (au *AuthUsecase) ForgotPassword(ctx context.Context, req dto.ForgotPasswordRequest) error {

	user, err := au.authRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return errors.New("User tidak ditemukan")
	}

	token, err := jwt.GenerateResetPasswordToken(user.ID.String(), user.UpdatedAt.Unix())
	if err != nil {
		return errors.New("Gagal generate token")
	}

	fmt.Println("FORGET PASSWORD TOKEN: ", token)

	link := os.Getenv("FE_URL") + "/reset-password"

	return email.SendResetPasswordEmail(user.Email, link, token)
}

func (au *AuthUsecase) ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) error {

	claims, err := jwt.ValidateToken(req.Token, os.Getenv("RESET_PASSWORD_SECRET"))
	if err != nil {
		return errors.New("Token tidak valid atau sudah kedaluwarsa")
	}

	if claims.Type != "reset_password" {
		return errors.New("Tipe token tidak valid")
	}

	user, err := au.authRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return errors.New("User tidak ditemukan")
	}

	if claims.UpdatedAt < user.UpdatedAt.Unix() {
		return errors.New("Token tidak valid atau sudah kedaluwarsa")
	}

	hashedPassword, err := bcrypt.HashPassword(req.Password)
	if err != nil {
		return errors.New("Gagal hash password")
	}

	user.Password = hashedPassword
	user.RefreshToken = ""
	return au.authRepo.Update(ctx, user)
}

func (au *AuthUsecase) Refresh(ctx context.Context, refreshToken string) (*dto.RefreshResponse, error) {

	claims, err := jwt.ValidateToken(refreshToken, os.Getenv("REFRESH_TOKEN_SECRET"))
	if err != nil {
		return nil, errors.New("Refresh token tidak valid atau kedaluwarsa")
	}

	if claims.Type != "refresh" {
		return nil, errors.New("Tipe token tidak valid")
	}

	user, err := au.authRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, errors.New("User tidak ditemukan")
	}

	if user.RefreshToken != refreshToken {
		return nil, errors.New("Refresh token tidak sesuai")
	}

	if claims.UpdatedAt < user.UpdatedAt.Unix() {
		return nil, errors.New("Refresh token sudah tidak valid, silahkan login kembali")
	}

	accessToken, err := jwt.GenerateAccessToken(user.ID.String(), string(user.Role), user.UpdatedAt.Unix())
	if err != nil {
		return nil, errors.New("Gagal generate access token")
	}

	return &dto.RefreshResponse{
		AccessToken: accessToken,
	}, nil
}

func (au *AuthUsecase) Logout(ctx context.Context, refreshToken string) error {
	claims, err := jwt.ValidateToken(refreshToken, os.Getenv("REFRESH_TOKEN_SECRET"))
	if err != nil {
		return err
	}

	user, err := au.authRepo.FindByID(ctx, claims.UserID)
	if err != nil || user == nil {
		return nil
	}

	user.RefreshToken = ""
	return au.authRepo.Update(ctx, user)
}
