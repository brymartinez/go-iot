package main

import (
	"go-iot/api/handler"
	"go-iot/api/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err.Error())
	}
	gin.DisableConsoleColor()
	router := gin.New()
	router.Use(gin.Recovery())
	router.GET("/device/:id", handler.GetDevice)
	router.POST("/device", handler.CreateDevice)
	router.PUT("/device/:id", handler.ConfigureDevice)
	router.PUT("/device/:id/activate", handler.ActivateDevice)
	// router.PUT("/device/:id/deactivate", handler.DeactivateDevice)
	// router.DELETE("/device/:id", handler.DeleteDevice)
	go service.Subscribe("localhost:3000/sns-endpoint")
	router.Run("localhost:3000")

}
