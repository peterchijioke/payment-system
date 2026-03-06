package handlers

import (
	"net/http"
	"strconv"

	"take-Home-assignment/internal/dto"
	"take-Home-assignment/internal/services"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentService *services.PaymentService
}

func NewPaymentHandler(paymentService *services.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

func (h *PaymentHandler) InitiatePayment(c *gin.Context) {
	idempotencyKey := c.GetHeader("Idempotency-Key")
	if idempotencyKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Idempotency-Key header is required",
		})
		return
	}

	var req dto.PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	resp, err := h.paymentService.ProcessPayment(h.paymentService.GetDB(), &req, idempotencyKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    resp,
	})
}

func (h *PaymentHandler) GetPayment(c *gin.Context) {
	paymentID := c.Param("id")

	payment, err := h.paymentService.GetTransactionDetails(paymentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "payment not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    payment,
	})
}

func (h *PaymentHandler) ListPayments(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	status := c.Query("status")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	payments, total, err := h.paymentService.ListPayments(limit, offset, status, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"payments": payments,
			"total":    total,
			"limit":    limit,
			"offset":   offset,
		},
	})
}

func (h *PaymentHandler) HandleWebhook(c *gin.Context) {
	signature := c.GetHeader("X-Webhook-Signature")
	rawBody := c.Request.Body

	payload, err := h.paymentService.ProcessWebhook(h.paymentService.GetDB(), rawBody, signature)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "webhook processed",
		"data":    payload,
	})
}

func (h *PaymentHandler) ListAccounts(c *gin.Context) {
	accounts, err := h.paymentService.ListAccounts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    accounts,
	})
}
