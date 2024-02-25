package main

import (
	"go-iot/api/handler"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err.Error())
	}
	router := gin.Default()
	router.GET("/device/:id", handler.GetDevice)
	router.POST("/device", handler.CreateDevice)
	// router.PUT("/payment/:id", handler.ConfigureDevice)
	// router.PUT("/payment/:id/activate", handler.ActivateDevice)
	// router.PUT("/payment/:id/deactivate", handler.ActivateDevice)
	// router.DELETE("/payment/:id", handler.DeleteDevice)
	router.Run("localhost:3000")
}
