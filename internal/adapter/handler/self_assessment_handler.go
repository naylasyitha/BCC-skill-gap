package handler

import (
	"net/http"
	"project-bcc/dto"
	"project-bcc/internal/usecase"
	"strings"

	"github.com/gin-gonic/gin"
)

type SelfAssessmentHandler struct {
	selfAssessmentUsecase *usecase.SelfAssessmentUsecase
}

func NewSelfAssessmentHandler(us *usecase.SelfAssessmentUsecase) *SelfAssessmentHandler {
	return &SelfAssessmentHandler{selfAssessmentUsecase: us}
}

func (s *SelfAssessmentHandler) SubmitAssessment(c *gin.Context) {
	careersessionID := c.Param("careerSessionId")
	if careersessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Career session ID tidak ada",
		})
		return
	}

	var req dto.SelfAssessmentRequest
	err := c.ShouldBindJSON(&req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	res, err := s.selfAssessmentUsecase.ProcessSelfAssessment(c.Request.Context(), careersessionID, req)
	if err != nil {
		if strings.Contains(err.Error(), "Status career session tidak sesuai") {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Terjadi error pada internal server",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Self Assessment berhasil disubmit",
		"data":    res,
	})

}
