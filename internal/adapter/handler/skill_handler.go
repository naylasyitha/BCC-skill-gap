package handler

import (
	"errors"
	"net/http"
	"project-bcc/dto"
	"project-bcc/internal/usecase"
	"strings"

	"github.com/gin-gonic/gin"
)

type SkillHandler struct {
	skillUsecase *usecase.SkillUsecase
}

func NewSkillHandler(s *usecase.SkillUsecase) *SkillHandler {
	return &SkillHandler{skillUsecase: s}
}

func (s *SkillHandler) GetAllSkill(c *gin.Context) {
	result, err := s.skillUsecase.GetAllSkill(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil daftar Skill",
		"data":    result,
	})
}

func (s *SkillHandler) GetSkillById(c *gin.Context) {
	result, err := s.skillUsecase.GetSkillById(c.Request.Context(), c.Param("skillId"))
	if err != nil {
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

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil detail Skill",
		"data":    result,
	})
}

func (s *SkillHandler) CreateSkill(c *gin.Context) {
	var req dto.SkillCreateRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	result, err := s.skillUsecase.CreateSkill(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Skill berhasil ditambahkan",
		"data":    result,
	})
}

func (s *SkillHandler) UpdateSkill(c *gin.Context) {
	var req dto.SkillUpdateRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	result, err := s.skillUsecase.UpdateSkill(c.Request.Context(), c.Param("skillId"), req)
	if err != nil {
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

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Skill berhasil diubah",
		"data":    result,
	})
}

func (s *SkillHandler) DeleteSkill(c *gin.Context) {
	err := s.skillUsecase.DeleteSkill(c.Request.Context(), c.Param("skillId"))
	if err != nil {
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

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Skill berhasil dihapus",
	})
}

func (s *SkillHandler) CareerSkillAsign(c *gin.Context) {
	var req dto.CareerSkillCreateRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	result, err := s.skillUsecase.CareerSkillAsign(c.Request.Context(), req)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "tidak valid") {
			c.JSON(http.StatusBadRequest, gin.H{
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
		"message": "Berhasil menambahkan skill ke karir",
		"data":    result,
	})
}

func (s *SkillHandler) UpdateCareerSkill(c *gin.Context) {
	var req dto.CareerSkillUpdateRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	result, err := s.skillUsecase.UpdateCareerSkill(c.Request.Context(), c.Param("skillId"), req)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "tidak ditemukan") {
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
		"message": "Skill berhasil diubah",
		"data":    result,
	})
}

func (s *SkillHandler) RemoveSkillFromCareer(c *gin.Context) {
	err := s.skillUsecase.RemoveSkillFromCareer(c.Request.Context(), c.Param("skillId"))
	if err != nil {
		// Menangkap pesan error custom dari usecase untuk not found
		if strings.Contains(strings.ToLower(err.Error()), "tidak ditemukan") {
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
		"message": "Skill berhasil dihapus dari karir",
	})
}
