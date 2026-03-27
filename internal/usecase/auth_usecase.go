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
		return errors.New("gagal mengirim email verifikasi" + err.Error())
	}

	return nil
}

func (au *AuthUsecase) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {

	user, err := au.authRepo.FindByEmail(ctx, req.Email)
	if err != nil || user == nil {
		return nil, errors.New("email atau password salah")
	}

	if !user.IsVerified {
		return nil, errors.New("akun belum diverifikasi, silakan cek email anda")
	}

	if !bcrypt.CheckPassword(user.Password, req.Password) {
		return nil, errors.New("email atau password salah")
	}

	accessToken, err := jwt.GenerateAccessToken(user.ID.String(), string(user.Role))
	if err != nil {
		return nil, errors.New("gagal generate access token")
	}

	refreshToken, err := jwt.GenerateRefreshToken(user.ID.String(), string(user.Role), req.RememberMe)
	if err != nil {
		return nil, errors.New("gagal generate refresh token")
	}

	fmt.Println("ACCESS TOKEN: ", accessToken)
	user.RefreshToken = refreshToken
	err = au.authRepo.Update(ctx, user)
	if err != nil {
		return nil, errors.New("gagal menyimpan refresh token")
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
		return errors.New("token verifikasi tidak valid atau sudah kedaluwarsa")
	}

	if claims.Type != "email_verification" {
		return errors.New("tipe token tidak valid")
	}

	user, err := au.authRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return errors.New("user tidak ditemukan")
	}

	if user.IsVerified {
		return errors.New("email sudah diverifikasi")
	}

	user.IsVerified = true
	return au.authRepo.Update(ctx, user)
}

func (au *AuthUsecase) ResendVerification(ctx context.Context, req dto.ResendVerificationRequest) error {

	user, err := au.authRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return errors.New("user tidak ditemukan")
	}

	if user.IsVerified {
		return errors.New("email sudah diverifikasi")
	}

	token, err := jwt.GenerateEmailVerificationToken(user.ID.String())
	if err != nil {
		return errors.New("gagal generate token")
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
		return errors.New("user tidak ditemukan")
	}

	token, err := jwt.GenerateResetPasswordToken(user.ID.String())
	if err != nil {
		return errors.New("gagal generate token")
	}

	fmt.Println("FORGET PASSWORD TOKEN: ", token)

	link := req.CallbackUrl
	if link == "" {
		link = os.Getenv("FE_URL") + "/reset-password"
	}

	return email.SendResetPasswordEmail(user.Email, link, token)
}

func (au *AuthUsecase) ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) error {

	claims, err := jwt.ValidateToken(req.Token, os.Getenv("RESET_PASSWORD_SECRET"))
	if err != nil {
		return errors.New("token tidak valid atau sudah kedaluwarsa")
	}

	if claims.Type != "reset_password" {
		return errors.New("tipe token tidak valid")
	}

	user, err := au.authRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return errors.New("user tidak ditemukan")
	}

	hashedPassword, err := bcrypt.HashPassword(req.Password)
	if err != nil {
		return errors.New("gagal hash password")
	}

	user.Password = hashedPassword
	return au.authRepo.Update(ctx, user)
}

func (au *AuthUsecase) Refresh(ctx context.Context, refreshToken string) (*dto.RefreshResponse, error) {

	claims, err := jwt.ValidateToken(refreshToken, os.Getenv("REFRESH_TOKEN_SECRET"))
	if err != nil {
		return nil, errors.New("refresh token tidak valid atau kedaluwarsa")
	}

	if claims.Type != "refresh" {
		return nil, errors.New("tipe token tidak valid")
	}

	user, err := au.authRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}

	if user.RefreshToken != refreshToken {
		return nil, errors.New("refresh token tidak sesuai")
	}

	accessToken, err := jwt.GenerateAccessToken(user.ID.String(), string(user.Role))
	if err != nil {
		return nil, errors.New("gagal generate access token")
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
