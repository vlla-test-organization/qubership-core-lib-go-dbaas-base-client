package model

import "fmt"

type DbaaSCreateDbError struct {
	HttpCode int
	Message  string
	Errors   error
}

func (d DbaaSCreateDbError) Error() string {
	return fmt.Sprintf("Response code: %d, Message: %s", d.HttpCode, d.Message)
}
