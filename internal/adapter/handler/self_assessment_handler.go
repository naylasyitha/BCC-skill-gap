package handler

import (
	"net/http"
	"project-bcc/dto"
	"project-bcc/internal/usecase"

	"github.com/gin-gonic/gin"
)

type SelfAssessmentHandler struct {
	selfAssessmentUsecase *usecase.SelfAssessmentUsecase
}

func NewSelfAssessmentHandler(us *usecase.SelfAssessmentUsecase) *SelfAssessmentHandler {
	return &SelfAssessmentHandler{selfAssessmentUsecase: us}
}

func (s *SelfAssessmentHandler) SubmitAssessment(c *gin.Context) {
	var req dto.SelfAssessmentRequest
	err := c.ShouldBindJSON(&req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User tidak ditemukan",
		})
		return
	}

	res, err := s.selfAssessmentUsecase.ProcessSelfAssessment(c.Request.Context(), userID.(string), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Self Assessment berhasil disubmit",
		"data":    res,
	})

}
