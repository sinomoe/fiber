package dto

type SendMessageRequest struct {
	To      string `json:"to" binding:"required"`
	Message string `json:"message" binding:"required"`
}

type SendMessageResponse struct {
}
