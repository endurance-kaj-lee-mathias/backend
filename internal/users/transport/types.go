package transport

import (
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/application"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/domain"
)

type Handler struct {
	service application.Service
}

type addSupportRequest struct {
	SupportID string `json:"supportId"`
}

type createUserRequest struct {
	Email string        `json:"email"`
	Roles []domain.Role `json:"roles"`
}

func NewHandler(s application.Service) *Handler {
	return &Handler{service: s}
}
