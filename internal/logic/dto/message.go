package dto

type SendMessageRequest struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Message string `json:"message"`
}

type SendMessageResponse struct {
}
