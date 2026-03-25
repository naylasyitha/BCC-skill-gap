package dto

type RegisterRequest struct {
	Fullname    string `json:"fullname" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8"`
	CallbackUrl string `json:"callbackUrl"`
}

type LoginRequest struct {
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required"`
	RememberMe bool   `json:"rememberMe"`
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
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type AuthResponse struct {
	AccessToken  string   `json:"accessToken"`
	RefreshToken string   `json:"-"`
	User         UserData `json:"user"`
}
