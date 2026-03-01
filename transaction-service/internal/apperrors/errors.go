package apperrors

import "errors"

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrDuplicateTransaction = errors.New("transaction already processed")
	ErrInsufficientFunds    = errors.New("insufficient funds")
)