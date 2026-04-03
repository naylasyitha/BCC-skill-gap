package handler

import (
	"net/http"
	"project-bcc/dto"
	"project-bcc/internal/usecase"
	"strings"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentUsecase *usecase.PaymentUsecase
}

func NewPaymentHandler(u *usecase.PaymentUsecase) *PaymentHandler {
	return &PaymentHandler{paymentUsecase: u}
}

func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User tidak ditemukan atau sesi telah habis",
		})
		return
	}
	res, err := h.paymentUsecase.CreatePayment(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": usecase.ErrInternalServer.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Token pembayaran berhasil dibuat",
		"data":    res,
	})
}

func (h *PaymentHandler) HandleNotification(c *gin.Context) {
	var req dto.MidtransNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Payload tidak valid",
		})
		return
	}

	err := h.paymentUsecase.HandleNotification(c.Request.Context(), req)
	if err != nil {
		if strings.Contains(err.Error(), "tidak ditemukan") {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": usecase.ErrInternalServer.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Notifikasi berhasil diproses",
	})
}

func (h *PaymentHandler) GetPaymentStatus(c *gin.Context) {
	orderID := c.Param("orderId")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "orderId tidak valid atau kosong",
		})
		return
	}

	tx, err := h.paymentUsecase.GetPaymentStatus(c.Request.Context(), orderID)
	if err != nil {
		if strings.Contains(err.Error(), "tidak ditemukan") {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": usecase.ErrInternalServer.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"order_id": tx.ID,
			"status":   tx.Status,
		},
	})
}
