package order

import (
	"mime/multipart"

	"github.com/prakasa1904/panji-express/internal/model"
	"github.com/prakasa1904/panji-express/internal/services/member"
	"gorm.io/datatypes"
)

// payload.MetaFile
type CreatePayload struct {
	model.Request
	Description string               `json:"description" validate:"max=256"`
	MetaFile    multipart.FileHeader `json:"meta_file"`
	Status      int                  `json:"status"`
	MemberID    uint64               `json:"member_id" validate:"required"`
}

type UpdatePayload struct {
	model.Request
	ID          uint64               `json:"id" validate:"required"`
	InvoiceID   string               `json:"invoice_id"` // generate once in step pickup by courier
	Description string               `json:"description" validate:"max=256"`
	MetaFile    multipart.FileHeader `json:"meta_file"`
	Status      int                  `json:"status"`
	CreatedAt   string               `json:"created_at" validate:"required"`
	MemberID    uint64               `json:"member_id" validate:"required"`
}

type DeletePayload struct {
	model.Request
	ID string `json:"id,omitempty"`
}

type Response struct {
	ID          uint64                                   `json:"id,omitempty"`
	InvoiceID   string                                   `json:"invoice_id,omitempty"`
	Description string                                   `json:"description,omitempty"`
	MetaFile    datatypes.JSONType[multipart.FileHeader] `json:"meta_file,omitempty"`
	Status      int                                      `json:"status,omitempty"`
	MemberID    uint64                                   `json:"member_id,omitempty"`
	CreatedAt   string                                   `json:"created_at,omitempty"`
	UpdatedAt   string                                   `json:"updated_at,omitempty"`
	Member      member.Response                          `json:"member,omitempty"`
}
