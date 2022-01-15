package comet

type SendRequest struct {
	From    string
	To      string
	Message string
}

type SendResponse struct{}
