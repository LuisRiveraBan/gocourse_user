package user

import (
	"errors"
	"fmt"
)

var ErrFirstNameRequired = errors.New("first name is required")
var ErrLastNameRequired = errors.New("last name is required")

// var ErrUserNotFound = errors.New("user not found")

type ErrNotFound struct {
	UserId string
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("user '%s doesnt't exist", e.UserId)
}
