package order

import (
	"time"

	"github.com/devetek/golang-webapp-boilerplate/internal/model"
	"gorm.io/datatypes"
)

type Response struct {
	ID          uint64                  `json:"id,omitempty"`
	InvoiceID   string                  `json:"invoice_id,omitempty"`
	Description string                  `json:"description,omitempty"`
	MetaFile    datatypes.JSONType[any] `json:"meta_file,omitempty"`
	Status      int                     `json:"status,omitempty"`
	CreatedAt   time.Time               `json:"created_at,omitempty"`
	UpdatedAt   time.Time               `json:"updated_at,omitempty"`
}

type CreateRequest struct {
	model.Request
	Description string                  `json:"description" validate:"required,max=256"`
	MetaFile    datatypes.JSONType[any] `json:"meta_file" validate:"required,max=256"`
}
