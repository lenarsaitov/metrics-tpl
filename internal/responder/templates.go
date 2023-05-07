package responder

const (
	defaultInternalErrorResponseMessage = "something going wrong"

	requestTemplate = `
{
	"response": {
		"text": "%s"
	}
}`
)
