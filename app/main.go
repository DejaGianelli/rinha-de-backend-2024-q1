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
	var err error
	connStr := "postgres://admin:123@localhost/rinha?sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected!")

	router := gin.Default()
	router.POST("/clientes/:id/transacoes", doTransactionHandler)
	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func doTransactionHandler(c *gin.Context) {
	customerId, err := strconv.Atoi(c.Param("id"))
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

	var customer customer
	row := db.QueryRow("SELECT limite, saldo_inicial FROM clientes WHERE id = $1", customerId)
	if err := row.Scan(&customer.Limit, &customer.Balance); err != nil {
		if err == sql.ErrNoRows {
			c.Status(http.StatusNotFound)
			return
		}
		log.Fatal(err)
	}

	if newTransaction.Type == "d" {
		newBalance := customer.Balance - newTransaction.Value
		if newBalance < (customer.Limit * -1) {
			c.Status(http.StatusUnprocessableEntity)
			return
		}
	}

	c.IndentedJSON(http.StatusOK, gin.H{"limite": customer.Limit, "saldo": customer.Balance})
}
