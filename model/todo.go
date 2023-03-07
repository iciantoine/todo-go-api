package model

import (
	"time"

	"github.com/google/uuid"
)

// Todo is the model.
type Todo struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	IsDone    bool      `json:"is_done"`
	Message   string    `json:"message" binding:"required"`
}
