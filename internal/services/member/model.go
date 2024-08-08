package member

import (
	"time"

	"github.com/prakasa1904/panji-express/internal/model"
	"github.com/prakasa1904/panji-express/internal/services/group"
)

type Response struct {
	ID              int64          `json:"id,omitempty"`
	Fullname        string         `json:"fullname,omitempty"`
	Email           string         `json:"email,omitempty"`
	Username        string         `json:"username,omitempty"`
	Password        string         `json:"password,omitempty"`
	ConfirmPassword string         `json:"confirm_password,omitempty"`
	CreatedAt       time.Time      `json:"created_at,omitempty"`
	UpdatedAt       time.Time      `json:"updated_at,omitempty"`
	GroupID         int64          `json:"group_id,omitempty"`
	Group           group.Response `json:"group,omitempty"`
}

type RequestPayload struct {
	model.Request
	ID              string `json:"id"`
	Fullname        string `json:"fullname"`
	Username        string `json:"username" validate:"max=100"`
	Email           string `json:"email" validate:"required,max=100"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	GroupID         string `json:"group_id" validate:"required"`
}

type DeletePayload struct {
	model.Request
	ID string `json:"id,omitempty"`
}

// TODO: Might be deprecate
type FindRequest struct {
	model.Request
	Fullname string `json:"fullname"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
