package handlers

import (
	"averroes-task/models"
	"averroes-task/store"
	"fmt"
	"net/http"
	"time"

	"github.com/blockloop/scan"
	"github.com/gin-gonic/gin"
)

func AddReview(c *gin.Context) {
	var reviewRequest models.ReviewRequest

	userData := getUserDataFromToken(c)

	dbClient := store.GetPostgres()

	if err := c.BindJSON(&reviewRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	user := GetUserByEmail(userData.Email)

	isWatched := isWatched(user.Id, reviewRequest.MovieId)

	if isWatched != true {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Movie must be added to watchlist"})
		c.Abort()
		return
	}

	insertQuery := `INSERT INTO reviews ("user_id", "movie_id", "review", "rate", "date") VALUES ($1, $2, $3, $4, $5)`

	_, insertErr := dbClient.Exec(insertQuery, user.Id, reviewRequest.MovieId, reviewRequest.Review, reviewRequest.Rate, time.Now())

	if insertErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": insertErr.Error()})
		c.Abort()
		return
	}

	numberOfRaters := getMovieReviewsCount(reviewRequest.MovieId, c)

	fmt.Println(numberOfRaters)

	// to add the current rate to the counter
	numberOfRaters.Count += 1
	numberOfRaters.RatesSum += reviewRequest.Rate

	updatedRate := float64(numberOfRaters.RatesSum) / float64(numberOfRaters.Count)

	fmt.Println(updatedRate)

	updateRateErr := updateRateById(reviewRequest.MovieId, updatedRate)

	if updateRateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": updateRateErr.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"created": true,
	})

}

func getMovieReviewsCount(movieId int, c *gin.Context) models.ReviewsCountAndRaters {
	var countAndSum models.ReviewsCountAndRaters

	dbClient := store.GetPostgres()

	query := `SELECT count(*), SUM(rate) as total FROM reviews WHERE movie_id = $1`

	rows, queryErr := dbClient.Query(query, movieId)

	if queryErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": queryErr.Error()})
		c.Abort()
	}

	scan.Row(&countAndSum, rows)

	return countAndSum
}
