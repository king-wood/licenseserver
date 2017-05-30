package internalerrors

type ErrorType int

const (
	_dummy ErrorType = iota
	InternalError
	RequestError
	AuthorizedError
	UnKnownError
)

type LogicError struct {
	Type        ErrorType
	Description string
}

func NewLogicError(t ErrorType, desc string) *LogicError {
	return &LogicError{t, desc}
}
