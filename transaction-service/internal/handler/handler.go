package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/grayDorian1/Entain/internal/apperrors"
	"github.com/grayDorian1/Entain/internal/logger"
	"github.com/grayDorian1/Entain/internal/model"
)

type Service interface {
	ProcessTransaction(ctx context.Context, userID uint64, sourceType, state, amount, transactionID string) error
	GetBalance(ctx context.Context, userID uint64) (string, error)
}

type Handler struct {
	svc Service
}

func New(svc Service) *Handler {
	return &Handler{svc: svc}
}

// HandleTransaction godoc
// @Summary      Apply transaction
// @Description  Increases or decreases user balance
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        userId          path      int                      true  "User ID"
// @Param        Source-Type     header    string                   true  "Source type" Enums(game, server, payment)
// @Param        request         body      model.TransactionRequest true  "Transaction payload"
// @Success      200
// @Failure      400  {object}  model.ErrorResponse
// @Failure      404  {object}  model.ErrorResponse
// @Failure      422  {object}  model.ErrorResponse
// @Router       /user/{userId}/transaction [post]
func (h *Handler) HandleTransaction(w http.ResponseWriter, r *http.Request) {
	userID, err := parseUserID(r)
	if err != nil {
		logger.Error(fmt.Sprintf("invalid userId: %v", err))
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	logger.Info(fmt.Sprintf("POST /user/%d/transaction source=%s", userID, r.Header.Get("Source-Type")))

	sourceType := r.Header.Get("Source-Type")

	var req model.TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error(fmt.Sprintf("invalid request body for user %d: %v", userID, err))
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	err = h.svc.ProcessTransaction(r.Context(), userID, sourceType, req.State, req.Amount, req.TransactionID)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrDuplicateTransaction):
			logger.Info(fmt.Sprintf("duplicate transaction %s for user %d", req.TransactionID, userID))
			w.WriteHeader(http.StatusOK)
			return
		case errors.Is(err, apperrors.ErrUserNotFound):
			logger.Error(fmt.Sprintf("user %d not found", userID))
			writeError(w, http.StatusNotFound, "user not found")
		case errors.Is(err, apperrors.ErrInsufficientFunds):
			logger.Error(fmt.Sprintf("insufficient funds for user %d", userID))
			writeError(w, http.StatusUnprocessableEntity, "insufficient funds")
		default:
			logger.Error(fmt.Sprintf("transaction failed for user %d: %v", userID, err))
			writeError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	logger.Info(fmt.Sprintf("transaction %s applied for user %d", req.TransactionID, userID))
	w.WriteHeader(http.StatusOK)
}

// HandleBalance godoc
// @Summary      Get user balance
// @Description  Returns current balance for a user
// @Tags         balance
// @Produce      json
// @Param        userId  path      int  true  "User ID"
// @Success      200     {object}  model.BalanceResponse
// @Failure      404     {object}  model.ErrorResponse
// @Router       /user/{userId}/balance [get]
func (h *Handler) HandleBalance(w http.ResponseWriter, r *http.Request) {
	userID, err := parseUserID(r)
	if err != nil {
		logger.Error(fmt.Sprintf("invalid userId: %v", err))
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	logger.Info(fmt.Sprintf("GET /user/%d/balance", userID))

	balance, err := h.svc.GetBalance(r.Context(), userID)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			logger.Error(fmt.Sprintf("user %d not found", userID))
			writeError(w, http.StatusNotFound, "user not found")
			return
		}
		logger.Error(fmt.Sprintf("get balance failed for user %d: %v", userID, err))
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	logger.Info(fmt.Sprintf("balance for user %d: %s", userID, balance))
	writeJSON(w, http.StatusOK, model.BalanceResponse{
		UserID:  userID,
		Balance: balance,
	})
}

func parseUserID(r *http.Request) (uint64, error) {
	param := r.PathValue("userId")
	id, err := strconv.ParseUint(param, 10, 64)
	if err != nil || id == 0 {
		return 0, errors.New("userId must be a positive integer")
	}
	return id, nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, model.ErrorResponse{Error: msg})
}