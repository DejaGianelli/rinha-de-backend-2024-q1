package main

import (
	"net/http"
	"slices"
	"strconv"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
)

type transaction struct {
	Type        string `json:"tipo"`
	Value       int    `json:"valor"`
	Description string `json:"descricao"`
}

func main() {
	router := gin.Default()
	router.POST("/clientes/:id/transacoes", doTransactionHandler)
	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func doTransactionHandler(c *gin.Context) {
	customerId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	var newTransaction transaction

	if err := c.BindJSON(&newTransaction); err != nil {
		return
	}

	if newTransaction.Value <= 0 {
		c.Status(http.StatusBadRequest)
		return
	}

	if !slices.Contains([]string{"c", "d"}, newTransaction.Type) {
		c.Status(http.StatusBadRequest)
		return
	}

	descSize := utf8.RuneCountInString(newTransaction.Description)
	if descSize < 1 || descSize > 10 {
		c.Status(http.StatusBadRequest)
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"id": customerId})
}
