package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"slices"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/DejaGianelli/rinha-de-backend-2024-q1/models"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type TransactionRequest struct {
	Type        string `json:"tipo"`
	Value       int    `json:"valor"`
	Description string `json:"descricao"`
}

type TransactionHandler struct {
	Db *sql.DB
}

func (handler *TransactionHandler) Handle(c *gin.Context) {
	customerId, err := strconv.Atoi(c.Param("id"))

	//Do Basic Validation
	if err != nil {
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	var newTransaction TransactionRequest

	if err := c.BindJSON(&newTransaction); err != nil {
		return
	}

	if newTransaction.Value <= 0 {
		c.Status(http.StatusUnprocessableEntity)
		return
	}
	if !slices.Contains([]string{"c", "d"}, newTransaction.Type) {
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	descSize := utf8.RuneCountInString(newTransaction.Description)
	if descSize < 1 || descSize > 10 {
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	//Initialize Database Transaction
	tx, err := handler.Db.BeginTx(c, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	//Do Bussiness Logic
	var customer models.Customer
	row := tx.QueryRowContext(c, "SELECT id, \"limit\", balance FROM customers WHERE id = $1 FOR UPDATE", customerId)
	if err := row.Scan(&customer.Id, &customer.Limit, &customer.Balance); err != nil {
		if err == sql.ErrNoRows {
			c.Status(http.StatusNotFound)
			return
		}
		log.Fatal(err)
	}

	var newBalance int
	if newTransaction.Type == "d" {
		newBalance = customer.Balance - newTransaction.Value
		if newBalance < (customer.Limit * -1) {
			c.Status(http.StatusUnprocessableEntity)
			return
		}
	} else {
		newBalance = customer.Balance + newTransaction.Value
	}

	_, err = tx.ExecContext(c, "INSERT INTO transactions (amount, type, customer_id, description, realized_at) VALUES ($1, $2, $3, $4, $5)",
		newTransaction.Value,
		newTransaction.Type,
		customer.Id,
		newTransaction.Description,
		time.Now())
	if err != nil {
		log.Fatal(err)
	}

	tx.ExecContext(c, "UPDATE customers SET balance = $1 WHERE id = $2", newBalance, customer.Id)

	customer.Balance = newBalance

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		log.Fatal(err)
	}

	c.IndentedJSON(http.StatusOK, gin.H{"limite": customer.Limit, "saldo": customer.Balance})
}
