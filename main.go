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
	router.PUT("/device/:id", handler.ConfigureDevice)
	// router.PUT("/device/:id/activate", handler.ActivateDevice)
	// router.PUT("/device/:id/deactivate", handler.DeactivateDevice)
	// router.DELETE("/device/:id", handler.DeleteDevice)
	router.Run("localhost:3000")
}
