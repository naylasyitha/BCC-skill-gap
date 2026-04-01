package handler

import (
	"net/http"
	"project-bcc/dto"
	"project-bcc/internal/usecase"

	"github.com/gin-gonic/gin"
)

type QuizHandler struct {
	quizUsecase usecase.QuizUsecase
}

func NewQuizHandler(qu usecase.QuizUsecase) *QuizHandler {
	return &QuizHandler{quizUsecase: qu}
}

func (h *QuizHandler) StartQuiz(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User tidak ditemukan atau sesi telah habis",
		})
		return
	}

	careerSessionID := c.Param("careerSessionId")
	if careerSessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "career session ID tidak valid",
		})
		return
	}

	res, err := h.quizUsecase.StartQuiz(c.Request.Context(), userID.(string), careerSessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Kuis berhasil dimulai",
		"data":    res,
	})
}

func (h *QuizHandler) UpdateAnswer(c *gin.Context) {
	quizSessionId := c.Param("quizSessionId")
	if quizSessionId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Quiz session ID tidak valid",
		})
		return
	}
	var req dto.UpdateAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Data tidak valid",
		})
		return
	}

	err := h.quizUsecase.UpdateAnswer(c.Request.Context(), quizSessionId, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Jawaban berhasil diperbarui",
	})
}

func (h *QuizHandler) SubmitQuiz(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User tidak ditemukan atau sesi telah habis",
		})
		return
	}

	quizSessionID := c.Param("quizSessionId")
	if quizSessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Quiz session ID tidak valid",
		})
		return
	}

	result, err := h.quizUsecase.SubmitQuiz(c.Request.Context(), userId.(string), quizSessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Kuis berhasil diserahkan",
		"data":    result,
	})
}
