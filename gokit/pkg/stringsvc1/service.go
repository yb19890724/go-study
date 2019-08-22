package stringsvc1

import (
	"errors"
	"strings"
)

var ErrEmpty = errors.New("empty string")

// StringService provides operations on strings.
type Serviceer interface {
	Uppercase(string) (string, error)
	Count(string) int
}

// stringService is a concrete implementation of StringService
type Service struct{}

func (Service) Uppercase(s string) (string, error) {
	if s == "" {
		return "", ErrEmpty
	}
	return strings.ToUpper(s), nil
}

func (Service) Count(s string) int {
	return len(s)
}
