package main

import (
	"go-iot/api/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/device/:id", handlers.GetDevice)
	// router.POST("/device", handlers.CreateDevice)
	// router.PUT("/payment/:id", handlers.ConfigureDevice)
	// router.PUT("/payment/:id/activate", handlers.ActivateDevice)
	// router.PUT("/payment/:id/deactivate", handlers.ActivateDevice)
	// router.DELETE("/payment/:id", handlers.DeleteDevice)
	router.Run("localhost:3000")
}
