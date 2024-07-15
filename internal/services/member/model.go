package member

import (
	"time"

	"github.com/devetek/golang-webapp-boilerplate/internal/model"
)

type Response struct {
	ID        uint64    `json:"id,omitempty"`
	Fullname  string    `json:"fullname,omitempty"`
	Email     string    `json:"email,omitempty"`
	Username  string    `json:"username,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type RegisterRequest struct {
	model.Request
	Fullname        string `json:"fullname"`
	Username        string `json:"username" validate:"required,max=100"`
	Email           string `json:"email" validate:"required,max=100"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
}
