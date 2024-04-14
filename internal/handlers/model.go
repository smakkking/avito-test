package handlers

import (
	"github.com/smakkking/avito_test/internal/services"
)

type Handler struct {
	bannerService *services.Service
}

func NewHandler(service *services.Service) *Handler {
	return &Handler{
		bannerService: service,
	}
}
