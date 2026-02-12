package transport

import "gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/message/application"

type Handler struct {
	service application.Service
}

func NewHandler(service application.Service) *Handler {
	return &Handler{service}
}
