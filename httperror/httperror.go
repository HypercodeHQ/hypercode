package httperror

import "fmt"

type HTTPError struct {
	StatusCode int
	Message    string
}

func (e HTTPError) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.StatusCode, e.Message)
}

func New(statusCode int, message string) error {
	return HTTPError{StatusCode: statusCode, Message: message}
}

func NotFound(message string) error {
	return HTTPError{StatusCode: 404, Message: message}
}

func BadRequest(message string) error {
	return HTTPError{StatusCode: 400, Message: message}
}

func Unauthorized(message string) error {
	return HTTPError{StatusCode: 401, Message: message}
}

func Forbidden(message string) error {
	return HTTPError{StatusCode: 403, Message: message}
}
