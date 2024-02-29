package main

import (
	"go-iot/api/handler"
	"go-iot/api/service"
	"log"
	"os"
	"os/signal"
	"syscall"

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
	// Start SNS subscription server
	go func() {
		err := service.Subscribe()
		if err != nil {
			log.Fatalf("Error starting SNS subscription server: %v", err)
		}
	}()

	// Start Gin HTTP server
	go func() {
		if err := router.Run("localhost:3000"); err != nil {
			log.Fatalf("Error starting Gin HTTP server: %v", err)
		}
	}()

	// Wait for termination signal
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down...")

	// Perform graceful shutdown (if needed)
	// Add cleanup logic here if necessary

	log.Println("Server gracefully stopped")

}
