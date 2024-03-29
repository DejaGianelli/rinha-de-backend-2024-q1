package main

import (
	"database/sql"
	"fmt"
	"github.com/DejaGianelli/rinha-de-backend-2024-q1/handlers"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func main() {
	//Load env variables
	DB_HOSTNAME := os.Getenv("DB_HOSTNAME")
	if DB_HOSTNAME == "" {
		DB_HOSTNAME = "localhost"
	}

	//Open Connection Pool to Postgres
	var db *sql.DB
	var err error

	connStr := "postgres://admin:123@" + DB_HOSTNAME + "/rinha?sslmode=require"

	db, err = sql.Open("postgres", connStr)
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(50)

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

	transactionHandler := handlers.TransactionHandler{
		Db: db,
	}

	extractHandler := handlers.ExtractHandler{
		Db: db,
	}

	router.POST("/clientes/:id/transacoes", transactionHandler.Handle)
	router.GET("/clientes/:id/extrato", extractHandler.Handler)

	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
