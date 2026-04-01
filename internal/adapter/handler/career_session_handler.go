package handler

import (
	"net/http"
	"project-bcc/dto"
	"project-bcc/internal/usecase"
	"strings"

	"github.com/gin-gonic/gin"
)

type CareerSessionHandler struct {
	careerSessionUsecase *usecase.CareerSessionUsecase
}

func NewCareerSessionHandler(cs *usecase.CareerSessionUsecase) *CareerSessionHandler {
	return &CareerSessionHandler{careerSessionUsecase: cs}
}

func (h *CareerSessionHandler) Create(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User tidak ditemukan atau sesi telah habis",
		})
		return
	}

	var req dto.CareerSessionCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	res, err := h.careerSessionUsecase.CreateCareerSession(c.Request.Context(), userID.(string), req)
	if err != nil {
		if strings.Contains(err.Error(), "mencapai batas") {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Terjadi masalah pada internal server",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Sesi karir berhasil dimulai",
		"data":    res,
	})
}
func (h *CareerSessionHandler) GetCareerSession(c *gin.Context) {
	careerSession, err := h.careerSessionUsecase.GetCareerSession(c.Request.Context(), c.Param("careerSessionId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil Mengambil Career Session",
		"data":    careerSession,
	})
}

func (h *CareerSessionHandler) GetAllCareerSession(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User tidak ditemukan atau sesi telah habis",
		})
		return
	}

	careerSessions, err := h.careerSessionUsecase.GetAllCareerSession(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Terjadi error pada internal server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil Mengambil List Career Session",
		"data":    careerSessions,
	})
}

func (h *CareerSessionHandler) GetDashboardAnalytics(c *gin.Context) {
	careerSessionID := c.Param("careerSessionId")
	if careerSessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Career session ID tidak valid",
		})
		return
	}

	analyticsData, err := h.careerSessionUsecase.GetDashboardAnalytics(c.Request.Context(), careerSessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Terjadi error pada internal server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil Mengambil Data Analytics",
		"data":    analyticsData,
	})
}
