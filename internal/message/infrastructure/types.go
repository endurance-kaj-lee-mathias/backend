package infrastructure

import (
	"database/sql"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/message/infrastructure/entities"
)

type Repository interface {
	ReadMessage() (entities.MessageEntity, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db}
}
