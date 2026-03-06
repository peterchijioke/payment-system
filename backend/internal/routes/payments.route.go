package routes

import (
	"take-Home-assignment/internal/handlers"

	"github.com/gin-gonic/gin"
)

type PaymentRoutes struct {
	handler *handlers.PaymentHandler
}

func NewPaymentRoutes(handler *handlers.PaymentHandler) *PaymentRoutes {
	return &PaymentRoutes{
		handler: handler,
	}
}

func (r *PaymentRoutes) RegisterRoutes(api *gin.RouterGroup) {
	payments := api.Group("/payments")
	{
		payments.POST("", r.handler.InitiatePayment)
		payments.GET("", r.handler.ListPayments)
		payments.GET("/:id", r.handler.GetPayment)
	}
	api.POST("/webhooks/provider", r.handler.HandleWebhook)
	api.GET("/accounts", r.handler.ListAccounts)
}
