package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"slices"
	"strconv"
	"unicode/utf8"
)

var db *sql.DB

type transaction struct {
	Type        string `json:"tipo"`
	Value       int    `json:"valor"`
	Description string `json:"descricao"`
}

type customer struct {
	Id      int
	Name    string
	Limit   int
	Balance int
}

func main() {
	//Open Connection Pool to Postgres
	var err error
	connStr := "postgres://admin:123@localhost/rinha?sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	//Ping Database do check connection
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected!")

	//Initialize Web Server
	router := gin.Default()
	router.POST("/clientes/:id/transacoes", doTransactionHandler)
	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func doTransactionHandler(c *gin.Context) {
	customerId, err := strconv.Atoi(c.Param("id"))

	//Do Basic Validation
	if err != nil {
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	var newTransaction transaction

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
	tx, err := db.BeginTx(c, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	//Do Bussiness Logic
	var customer customer
	row := tx.QueryRowContext(c, "SELECT id, limite, saldo_inicial FROM clientes WHERE id = $1 FOR UPDATE", customerId)
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

	_, err = tx.ExecContext(c, "INSERT INTO transacoes (amount, type, customer_id) VALUES ($1, $2, $3)", newTransaction.Value, newTransaction.Type, customer.Id)
	if err != nil {
		log.Fatal(err)
	}

	tx.ExecContext(c, "UPDATE clientes SET saldo_inicial = $1 WHERE id = $2", newBalance, customer.Id)

	customer.Balance = newBalance

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		log.Fatal(err)
	}

	c.IndentedJSON(http.StatusOK, gin.H{"limite": customer.Limit, "saldo": customer.Balance})
}
