package domain

import "errors"

var (
	ErrAlertNotFound = errors.New("alert not found")
	ErrAlertNotOpen  = errors.New("alert is not open")
)
