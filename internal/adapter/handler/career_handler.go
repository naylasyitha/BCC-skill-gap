package handler

import (
	"errors"
	"net/http"
	"project-bcc/dto"
	"project-bcc/internal/usecase"
	"strings"

	"github.com/gin-gonic/gin"
)

type CareerHandler struct {
	careerUsecase *usecase.CareerUsecase
}

func NewCareerHandler(car *usecase.CareerUsecase) *CareerHandler {
	return &CareerHandler{careerUsecase: car}
}

func (ch *CareerHandler) GetAllCareer(c *gin.Context) {
	result, err := ch.careerUsecase.GetAllCareer(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil daftar karir",
		"data":    result,
	})
}

func (ch *CareerHandler) GetCareerById(c *gin.Context) {
	result, err := ch.careerUsecase.GetCareerById(c.Request.Context(), c.Param("careerId"))
	if err != nil {
		if errors.Is(err, usecase.ErrCareerNotFound) {
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
		"message": "Berhasil mengambil detail Karir",
		"data":    result,
	})
}

func (ch *CareerHandler) CreateCareer(c *gin.Context) {
	var req dto.CareerCreateRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	result, err := ch.careerUsecase.CreateCareer(c.Request.Context(), req)
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
		"message": "Karir berhasil dibuat",
		"data":    result,
	})
}

func (h *CareerHandler) UpdateCareer(c *gin.Context) {
	careerID := c.Param("careerId")

	var req dto.CareerEditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	result, err := h.careerUsecase.UpdateCareer(c.Request.Context(), careerID, req)
	if err != nil {
		if errors.Is(err, usecase.ErrCareerNotFound) || errors.Is(err, usecase.ErrSkillNotFound) {
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
		"message": "Karir berhasil diubah",
		"data":    result,
	})
}

func (ch *CareerHandler) DeleteCareer(c *gin.Context) {
	err := ch.careerUsecase.DeleteCareer(c.Request.Context(), c.Param("careerId"))
	if err != nil {
		if errors.Is(err, usecase.ErrCareerNotFound) {
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
		"message": "Karir berhasil dihapus",
	})
}
