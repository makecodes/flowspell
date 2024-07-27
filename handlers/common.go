package handlers

import (
	"fmt"
)

type CustomError struct {
	Message map[string]string
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("%v", e.Message)
}
