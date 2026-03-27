package dto

type RegisterRequest struct {
	Fullname    string `json:"full_name" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8"`
	CallbackUrl string `json:"callbackUrl"`
}

type LoginRequest struct {
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required"`
	RememberMe bool   `json:"remember_me"`
}

type VerifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
}

type ResendVerificationRequest struct {
	Email       string `json:"email" binding:"required,email"`
	CallbackUrl string `json:"callbackUrl"`
}

type ForgotPasswordRequest struct {
	Email       string `json:"email" binding:"required,email"`
	CallbackUrl string `json:"callbackUrl"`
}

type ResetPasswordRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

type UserData struct {
	ID       string `json:"id"`
	Fullname string `json:"full_name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type AuthResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"-"`
	User         UserData `json:"user"`
}
type RefreshResponse struct {
	AccessToken string `json:"access_token"`
}
