package handlers

import "take-Home-assignment/internal/services"

type Container struct {
	Payment *PaymentHandler
}

func InitHandlers(services *services.Container) *Container {
	return &Container{
		Payment: NewPaymentHandler(services.Payment),
	}
}
