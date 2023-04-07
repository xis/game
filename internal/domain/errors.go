package domain

import "errors"

var (
	ErrResourceNotFound = errors.New("resource not found")
	ErrInternal         = errors.New("internal error")
)
