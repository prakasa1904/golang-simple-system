package group

import (
	"time"

	"github.com/devetek/golang-webapp-boilerplate/internal/model"
)

type Response struct {
	ID        uint64    `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type CreateRequest struct {
	model.Request
	Name string `json:"name" validate:"required,max=256"`
}
