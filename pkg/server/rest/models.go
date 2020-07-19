package rest

type (
	V1Base struct {
		ID      string `json:"id"`
		Created int64  `json:"created"`
		Message string `json:"message"`
		To      string `json:"to"`
	}

	ListV1Res struct {
		Reminders []string `json:"reminders"`
	}

	CreateAtV1Req struct {
		V1Base

		Moment int64 `json:"moment"`
	}

	CreateAfterV1Req struct {
		V1Base

		After int64 `json:"after"`
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
