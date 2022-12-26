package handlers

import (
	"averroes-task/models"
	"averroes-task/store"
	"net/http"
	"time"

	"github.com/blockloop/scan"
	"github.com/gin-gonic/gin"
)

func AddMovie(c *gin.Context) {
	var addRequest models.AddMovieRequest
	userData := getUserDataFromToken(c)

	dbClient := store.GetPostgres()

	if err := c.BindJSON(&addRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	user := GetUserByEmail(userData.Email)

	timeFormat, _ := time.Parse("2006/01/02", addRequest.Date)

	query := `INSERT INTO movies ("name", "description", "cover", "date", "rate", "user_created_id") VALUES ($1, $2, $3, $4, $5, $6)`

	_, insertErr := dbClient.Exec(query, addRequest.Name, addRequest.Description, addRequest.Cover, timeFormat, 0, user.Id)

	if insertErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": insertErr.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"created": true,
	})
}

func UpdateMovie(c *gin.Context) {
	var addRequest models.AddMovieRequest

	movieId := c.Param("id")

	dbClient := store.GetPostgres()

	if err := c.BindJSON(&addRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	timeFormat, _ := time.Parse("2006/01/02", addRequest.Date)

	updateQuery := `UPDATE movies SET name = $1, description = $2, cover = $3, date = $4 where id = $5`

	_, err := dbClient.Exec(updateQuery, addRequest.Name, addRequest.Description, addRequest.Cover, timeFormat, movieId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"updated": true,
	})
}

func DeleteMovie(c *gin.Context) {
	movieId := c.Param("id")

	dbClient := store.GetPostgres()

	deleteQuery := `DELETE FROM movies WHERE id = $1`

	_, err := dbClient.Exec(deleteQuery, movieId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"deleted": true,
	})
}

func GetMovieById(c *gin.Context) {
	movieId := c.Param("id")

	movie := getMovieById(movieId)

	c.JSON(http.StatusOK, movie)
}

func ListMovies(c *gin.Context) {
	var movies []models.Movie

	sortBy := c.Query("sortBy")
	dir := c.Query("dir")

	dbClient := store.GetPostgres()

	query := `SELECT * FROM movies`

	if sortBy != "" {
		query = `SELECT * FROM movies ORDER BY ` + sortBy

		if dir != "" {
			query = query + " " + dir
		}
	}

	rows, queryErr := dbClient.Query(query)

	if queryErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": queryErr.Error()})
		c.Abort()
		return
	}

	scan.Rows(&movies, rows)

	c.JSON(http.StatusOK, gin.H{
		"movies": movies,
	})
}

func getMovieById(id string) models.Movie {
	dbClient := store.GetPostgres()

	var movie models.Movie

	query := `SELECT * FROM users WHERE id = $1`

	rows, _ := dbClient.Query(query, id)

	scan.Row(&movie, rows)

	return movie
}

func updateRateById(movieId int, rate float64) error {
	dbClient := store.GetPostgres()

	updateQuery := `UPDATE movies SET rate = $1 WHERE id = $2`

	_, err := dbClient.Exec(updateQuery, rate, movieId)

	if err != nil {
		return err
	}

	return nil
}
