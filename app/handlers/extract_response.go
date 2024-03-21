package handlers

import (
	"time"
)

type ExtractResponse struct {
	Balance          BalanceResponse           `json:"saldo"`
	LastTransactions []LastTransactionResponse `json:"ultimas_transacoes"`
}

type BalanceResponse struct {
	Total       int       `json:"total"`
	ExtractDate time.Time `json:"extract_date"`
	Limit       int       `json:"limit"`
}

type LastTransactionResponse struct {
	Value       int       `json:"valor"`
	Type        string    `json:"tipo"`
	Description string    `json:"descricao"`
	RealizedAt  time.Time `json:"realizada_em"`
}
