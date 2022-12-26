package handlers

import (
	"averroes-task/models"
	"averroes-task/store"
	"net/http"
	"time"

	"github.com/blockloop/scan"
	"github.com/gin-gonic/gin"
)

func AddToWatchlist(c *gin.Context) {
	var addRequest models.AddToWatchlistRequest

	userData := getUserDataFromToken(c)

	dbClient := store.GetPostgres()

	if err := c.BindJSON(&addRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	user := GetUserByEmail(userData.Email)

	query := `INSERT INTO watchlist ("user_id", "movie_id", "date") VALUES ($1, $2, $3)`

	_, insertErr := dbClient.Exec(query, user.Id, addRequest.MovieId, time.Now())

	if insertErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": insertErr.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"created": true,
	})
}

func isWatched(userId int, movieId int) bool {
	var watched bool

	dbClient := store.GetPostgres()

	query := `SELECT EXISTS (SELECT id FROM watchlist WHERE user_id = $1 AND movie_id = $2)`

	rows, _ := dbClient.Query(query, userId, movieId)

	scan.Row(&watched, rows)

	return watched
}
