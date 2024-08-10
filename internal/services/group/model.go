package group

import (
	"time"

	"github.com/prakasa1904/panji-express/internal/model"
)

type Response struct {
	ID        int64     `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Status    int       `json:"status,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type RequestPayload struct {
	model.Request
	ID     uint64 `json:"id,omitempty"`
	Name   string `json:"name,omitempty" validate:"required,max=100"`
	Status int    `json:"status,omitempty" validate:"required"`
}

type DeletePayload struct {
	model.Request
	ID uint64 `json:"id,omitempty"`
}

type ResponseMutation struct {
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}
