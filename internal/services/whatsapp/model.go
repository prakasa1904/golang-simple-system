package whatsapp

import (
	"github.com/prakasa1904/panji-express/internal/model"
)

type SendMessagePayload struct {
	model.Request
	Receiver string `json:"receiver" validate:"required,max=100"`
	Message  string `json:"message" validate:"required,max=100"`
}
