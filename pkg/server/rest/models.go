package rest

type (
	V1 struct {
		ID      string `json:"id"`
		Created int64  `json:"created"`
		After   int64  `json:"after"`
		Message string `json:"message"`
		To      string `json:"to"`
	}

	Error struct {
		Message string `json:"message"`
		Err     string `json:"error"`
	}
)

func NewError(message string, err string) Error {
	return Error{
		Message: message,
		Err:     err,
	}
}

func (e *Error) Error() string {
	return e.Message
}
