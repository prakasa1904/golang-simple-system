package model

type Request struct {
	CurrentURL string `json:"current_url,omitempty"`
}

type ResponseMutation struct {
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}
