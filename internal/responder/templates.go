package responder

const (
	defaultInternalErrorResponseMessage = "something going wrong"
)

type textResponse struct {
	Text string `json:"text,omitempty"`
}

type serverResponse struct {
	Response textResponse `json:"response,omitempty"`
}
