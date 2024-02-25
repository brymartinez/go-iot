package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func InternalServerError(c *gin.Context) {
	c.IndentedJSON(http.StatusInternalServerError, gin.H{
		"error":   9999,
		"message": "Internal Server Error",
	})
}

func NotFoundError(c *gin.Context) {
	c.IndentedJSON(http.StatusNotFound, gin.H{
		"error":   1000,
		"message": "Not found",
	})
}
