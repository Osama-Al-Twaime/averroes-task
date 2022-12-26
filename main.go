package main

import (
	"averroes-task/auth"
	"averroes-task/handlers"
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
	r.POST("/users/register", handlers.RegisterNewUser)

	// login
	r.POST("/users/login", handlers.Login)

	// add movie to watchlist
	r.POST("/users/watchlist/add", auth.Auth(), handlers.AddToWatchlist)

	// get movie
	r.GET("/movies/:id", handlers.GetMovieById)

	// add movie
	r.POST("/movies/add", auth.Auth(), handlers.AddMovie)

	// update movie
	r.POST("/movies/edit/:id", auth.Auth(), handlers.UpdateMovie)

	// delete movie
	r.DELETE("/movies/delete/:id", auth.Auth(), handlers.DeleteMovie)

	// list movies
	r.GET("/movies/list", handlers.ListMovies)

	// rate and review movie
	r.POST("/movies/review", auth.Auth(), handlers.AddReview)

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

	// to handle request with go routines (concurrency)
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
