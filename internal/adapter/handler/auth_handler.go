package handler

import (
	"errors"
	"net/http"
	"os"
	"project-bcc/dto"
	"project-bcc/internal/usecase"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUsecase *usecase.AuthUsecase
}

func NewAuthHandler(uc *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: uc}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	//cek validasi request client
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	//melempar ke usecase
	err = h.authUsecase.Register(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	//mengirim respon
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Registrasi berhasil, silahkan cek email untuk verifikasi",
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	result, err := h.authUsecase.Login(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		if errors.Is(err, usecase.ErrUnverifiedAccount) {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// Set Durasi cookie
	maxAgeRefresh := 24 * 3600
	if req.RememberMe {
		maxAgeRefresh = 7 * 24 * 3600
	}

	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("refresh_token", result.RefreshToken, maxAgeRefresh, "/", "", os.Getenv("APP_ENV") == "production", true)
	c.SetCookie("role", result.User.Role, maxAgeRefresh, "/", "", os.Getenv("APP_ENV") == "production", true)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login berhasil",
		"data": gin.H{
			"access_token": result.AccessToken,
			"user":         result.User,
		},
	})
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var req dto.VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	err := h.authUsecase.VerifyEmail(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Email berhasil diverifikasi",
	})
}

func (h *AuthHandler) ResendVerification(c *gin.Context) {
	var req dto.ResendVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	err := h.authUsecase.ResendVerification(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Email verifikasi berhasil dikirim ulang",
	})
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req dto.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	err := h.authUsecase.ForgotPassword(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Link reset password telah dikirim ke email",
	})
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	err := h.authUsecase.ResetPassword(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Password berhasil direset",
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Refresh token tidak ditemukan",
		})
		return
	}

	result, err := h.authUsecase.Refresh(c.Request.Context(), refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Token berhasil diperbarui",
		"data":    result,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	refreshToken, _ := c.Cookie("refresh_token")

	if refreshToken != "" {
		h.authUsecase.Logout(c.Request.Context(), refreshToken)
	}

	// Hapus cookie
	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("refresh_token", "", -1, "/", "", os.Getenv("APP_ENV") == "production", true)
	c.SetCookie("role", "", -1, "/", "", os.Getenv("APP_ENV") == "production", false)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logout berhasil",
	})
}
