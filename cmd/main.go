package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"flashcard/internal/controller"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

func main() {
	log.Println("Starting the movie metadata service")
	f, err := os.Open("configs/base.yaml")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var cfg serviceConfig
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		panic(err)
	}
	log.Printf("Starting the metadata service on port %v", cfg.APIConfig.Port)

	r := gin.Default()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	// Listen for OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		fmt.Println("Server is shutting down...")

		// Create a context with a timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Gracefully shutdown the server
		if err := srv.Shutdown(ctx); err != nil {
			fmt.Printf("Server forced to shutdown: %v\n", err)
		}
	}()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/health", controller.Health)
	r.Run("localhost:" + cfg.APIConfig.Port)
}
