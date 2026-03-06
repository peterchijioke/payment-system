package routes

import (
	"take-Home-assignment/internal/handlers"
	"take-Home-assignment/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine, h *handlers.Container) {
	r.Use(middlewares.CORS())
	api := r.Group("/api")
	v1 := api.Group("/v1")
	paymentRoutes := NewPaymentRoutes(h.Payment)
	paymentRoutes.RegisterRoutes(v1)
}
