package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var validSourceTypes = map[string]struct{}{
	"game":    {},
	"server":  {},
	"payment": {},
}

var validStates = map[string]struct{}{
	"win":  {},
	"lose": {},
}

type TransactionRepository interface {
	ApplyTransaction(ctx context.Context, userID uint64, transactionID, sourceType, state, amount string) error
	GetBalance(ctx context.Context, userID uint64) (string, error)
}

type service struct {
	repo TransactionRepository
}

func New(repo TransactionRepository) *service {
	return &service{repo: repo}
}

func (s *service) ProcessTransaction(ctx context.Context, userID uint64, sourceType, state, amount, transactionID string) error {
	if err := validateSourceType(sourceType); err != nil {
		return err
	}
	if err := validateState(state); err != nil {
		return err
	}
	if err := validateAmount(amount); err != nil {
		return err
	}
	if strings.TrimSpace(transactionID) == "" {
		return fmt.Errorf("transactionId is required")
	}

	return s.repo.ApplyTransaction(ctx, userID, transactionID, sourceType, state, amount)
}

func (s *service) GetBalance(ctx context.Context, userID uint64) (string, error) {
	return s.repo.GetBalance(ctx, userID)
}

func validateSourceType(st string) error {
	if _, ok := validSourceTypes[st]; !ok {
		return fmt.Errorf("invalid source type: %q, must be one of: game, server, payment", st)
	}
	return nil
}

func validateState(state string) error {
	if _, ok := validStates[state]; !ok {
		return fmt.Errorf("invalid state: %q, must be win or lose", state)
	}
	return nil
}

func validateAmount(amount string) error {
	if strings.TrimSpace(amount) == "" {
		return fmt.Errorf("amount is required")
	}
	val, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return fmt.Errorf("amount must be a valid number")
	}
	if val <= 0 {
		return fmt.Errorf("amount must be positive")
	}
	// Check max 2 decimal places
	parts := strings.Split(amount, ".")
	if len(parts) == 2 && len(parts[1]) > 2 {
		return errors.New("amount must have at most 2 decimal places")
	}
	return nil
}