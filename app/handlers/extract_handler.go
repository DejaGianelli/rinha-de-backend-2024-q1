package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/DejaGianelli/rinha-de-backend-2024-q1/models"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type ExtractHandler struct {
	Db *sql.DB
}

func (handler *ExtractHandler) Handler(c *gin.Context) {
	customerId, err := strconv.Atoi(c.Param("id"))

	tx, err := handler.Db.BeginTx(c, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	var customer models.Customer
	row := tx.QueryRowContext(c, "SELECT id, balance, \"limit\" FROM customers WHERE id = $1", customerId)
	if err := row.Scan(&customer.Id, &customer.Balance, &customer.Limit); err != nil {
		if err == sql.ErrNoRows {
			c.Status(http.StatusNotFound)
			return
		}
	}

	var transactions []models.Transacao
	rows, err := tx.QueryContext(c, "SELECT amount, type, description, realized_at FROM transactions WHERE customer_id = $1 ORDER BY realized_at DESC LIMIT 10", customerId)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var transaction models.Transacao
		if err := rows.Scan(&transaction.Amount, &transaction.Type, &transaction.Description, &transaction.RealizedAt); err != nil {
			log.Fatal(err)
		}
		transactions = append(transactions, transaction)
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		log.Fatal(err)
	}

	lastTransactions := make([]LastTransactionResponse, 0)
	for _, t := range transactions {
		lastTransactions = append(lastTransactions, LastTransactionResponse{
			Value:       t.Amount,
			Type:        t.Type,
			Description: t.Description,
			RealizedAt:  t.RealizedAt,
		})
	}

	response := ExtractResponse{
		Balance: BalanceResponse{
			Total:       customer.Balance,
			ExtractDate: time.Now(),
			Limit:       customer.Limit,
		},
		LastTransactions: lastTransactions,
	}

	c.IndentedJSON(http.StatusOK, response)
}
