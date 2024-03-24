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
	ExtractDate time.Time `json:"data_extrato"`
	Limit       int       `json:"limite"`
}

type LastTransactionResponse struct {
	Value       int       `json:"valor"`
	Type        string    `json:"tipo"`
	Description string    `json:"descricao"`
	RealizedAt  time.Time `json:"realizada_em"`
}
