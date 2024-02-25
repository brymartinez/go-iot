package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func ConfigureDevice(c *gin.Context) {
	id := c.Param("id")
	fmt.Printf(id)

	// c.IndentedJSON(200, updatedDevice)

}
