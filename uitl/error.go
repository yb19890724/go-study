package uitl

type LimitError struct {
	Code    int
	Message string
}

func NewError (c int ,msg string) error {
	return &LimitError{Code:c ,Message:msg}
}

func (le LimitError) Error() string {
	return le.Message
}