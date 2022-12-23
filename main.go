package main

import (
	"averroes-task/handlers"
	"averroes-task/store"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// health check api
	r.GET("/status", healthCheck)

	// register a new user
	r.POST("/register", handlers.RegisterNewUser)

	// serving docs
	r.Static("/docs", "./docs")

	// initiating the server
	InitServer(r)
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}

func InitServer(r *gin.Engine) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	svr := &http.Server{
		Handler:      r,
		Addr:         ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println(store.GetPostgres())

	go func() {
		// start the web server on port and accept requests
		if err := svr.ListenAndServe(); err != nil {
			fmt.Printf("Error listening to server %s \n", err)
		}
	}()

	var wait time.Duration

	gracefulStop := make(chan os.Signal, 1)

	// We'll accept graceful shutdowns when quit via SIGINT, SIGTERM, SIGKILL, or Interrupt
	signal.Notify(gracefulStop, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)

	// Block until we receive our signal.
	<-gracefulStop

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	svr.Shutdown(ctx)

	fmt.Println("Server is shutting down")

	close(gracefulStop)

	os.Exit(0)
}
