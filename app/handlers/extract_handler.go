package handlers

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"net/http"

	"database/sql"
)

type ExtractHandler struct {
	Db *sql.DB
}

func (handler *ExtractHandler) Handler(c *gin.Context) {
	c.Status(http.StatusOK)
}
