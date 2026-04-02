package handler

import (
	"errors"
	"net/http"
	"project-bcc/dto"
	"project-bcc/internal/usecase"
	"strings"

	"github.com/gin-gonic/gin"
)

type QuestionHandler struct {
	questionUsecase *usecase.QuestionUsecase
}

func NewQuestionHandler(cs *usecase.QuestionUsecase) *QuestionHandler {
	return &QuestionHandler{questionUsecase: cs}
}

func (h *QuestionHandler) CreateQuestion(c *gin.Context) {
	var req dto.QuestionCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	res, err := h.questionUsecase.CreateQuestion(c.Request.Context(), req)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "tidak valid") {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}

		if errors.Is(err, usecase.ErrSkillNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
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

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Question berhasil dibuat",
		"data":    res,
	})
}

func (h *QuestionHandler) GetAllQuestion(c *gin.Context) {
	res, err := h.questionUsecase.GetAllQuestion(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil semua question",
		"data":    res,
	})
}

func (h *QuestionHandler) GetQuestionById(c *gin.Context) {
	id := c.Param("questionId")

	res, err := h.questionUsecase.GetQuestionById(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, usecase.ErrQuestionNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
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

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil detail question",
		"data":    res,
	})
}

func (h *QuestionHandler) UpdateQuestion(c *gin.Context) {
	id := c.Param("questionId")

	var req dto.QuestionUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	res, err := h.questionUsecase.UpdateQuestion(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, usecase.ErrQuestionNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
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

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Question berhasil diperbarui",
		"data":    res,
	})
}

func (h *QuestionHandler) DeleteQuestion(c *gin.Context) {
	id := c.Param("questionId")

	err := h.questionUsecase.DeleteQuestion(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, usecase.ErrQuestionNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
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

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Question berhasil dihapus",
	})
}
