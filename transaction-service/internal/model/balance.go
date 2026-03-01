package model 

type BalanceResponse struct {
	UserID  uint64 `json:"userId"`
	Balance string `json:"balance"`
}