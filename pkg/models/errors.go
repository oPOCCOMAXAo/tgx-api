package models

import "errors"

var (
	ErrInvalidParams = errors.New("invalid params")
	ErrNotFound      = errors.New("not found")
	ErrNotReady      = errors.New("not ready")
)
