package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/grayDorian1/Entain/internal/apperrors"
)

type repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *repository {
	return &repository{db: db}
}

func (r *repository) ApplyTransaction(ctx context.Context, userID uint64, transactionID, sourceType, state string, amount string) error {
	delta := amount
	if state == "lose" {
		delta = "-" + amount
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`INSERT INTO payments.transactions (transaction_id, user_id, source_type, state, amount)
		 VALUES ($1, $2, $3, $4, $5)`,
		transactionID, userID, sourceType, state, amount,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return apperrors.ErrDuplicateTransaction
		}
		return fmt.Errorf("insert transaction: %w", err)
	}

	tag, err := tx.Exec(ctx,
		`UPDATE accounts.users SET balance = balance + $1 WHERE id = $2 AND balance + $1 >= 0`,
		delta, userID,
	)
	if err != nil {
		return fmt.Errorf("update balance: %w", err)
	}

	if tag.RowsAffected() == 0 {
		var exists bool
		err = tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM accounts.users WHERE id = $1)`, userID).Scan(&exists)
		if err != nil {
			return fmt.Errorf("check user: %w", err)
		}
		if !exists {
			return apperrors.ErrUserNotFound
		}
		return apperrors.ErrInsufficientFunds
	}

	return tx.Commit(ctx)
}

func (r *repository) GetBalance(ctx context.Context, userID uint64) (string, error) {
	var balance string
	err := r.db.QueryRow(ctx,
		`SELECT TO_CHAR(balance, 'FM999999999999999990.00') FROM accounts.users WHERE id = $1`,
		userID,
	).Scan(&balance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", apperrors.ErrUserNotFound
		}
		return "", fmt.Errorf("query balance: %w", err)
	}
	return balance, nil
}