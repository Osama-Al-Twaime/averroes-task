package handlers

import (
	"averroes-task/auth"
	"averroes-task/models"
	"averroes-task/store"
	"net/http"
	"strings"

	"github.com/blockloop/scan"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func RegisterNewUser(c *gin.Context) {
	var userRequestBody models.UserRegisterRequest

	if err := c.BindJSON(&userRequestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	if hashErr := hashPassword(userRequestBody.Password, &userRequestBody); hashErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": hashErr.Error()})
		c.Abort()
		return
	}

	dbClient := store.GetPostgres()

	query := `INSERT INTO users ("name", "email", "password", "age") VALUES($1, $2, $3, $4)`

	_, insertErr := dbClient.Exec(query, userRequestBody.Name, userRequestBody.Email, userRequestBody.Password, userRequestBody.Age)

	if insertErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": insertErr.Error()})
		c.Abort()
		return
	}

	token, generateTokenErr := auth.GenerateJWT(userRequestBody.Name, userRequestBody.Email, userRequestBody.Age)

	if generateTokenErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": generateTokenErr.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"token": token,
	})
}

func Login(c *gin.Context) {
	var loginRequest models.LoginRequest

	var user models.User

	if err := c.BindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	dbClient := store.GetPostgres()

	query := `SELECT * FROM users WHERE email=$1`

	rows, queryErr := dbClient.Query(query, loginRequest.Email)

	if queryErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": queryErr.Error()})
		c.Abort()
		return
	}

	scan.Row(&user, rows)

	cehckPasswordErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
	if cehckPasswordErr != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "incorrect password"})
		c.Abort()
		return
	}

	token, generateTokenErr := auth.GenerateJWT(user.Name, user.Email, user.Age)

	if generateTokenErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": generateTokenErr.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"name":  user.Name,
		"email": user.Email,
		"age":   user.Age,
		"id":    user.Id,
	})
}

func hashPassword(password string, user *models.UserRegisterRequest) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

func GetUserByEmail(email string) models.User {
	dbClient := store.GetPostgres()

	var user models.User

	query := `SELECT * FROM users WHERE email = $1`

	rows, _ := dbClient.Query(query, email)

	scan.Row(&user, rows)

	return user
}

func getUserDataFromToken(c *gin.Context) *auth.JWTClaim {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(401, gin.H{"error": "request does not contain an authorization header"})
		c.Abort()
		return nil
	}

	authparts := strings.Split(authHeader, " ")

	if len(authparts) < 2 {
		c.JSON(401, gin.H{"error": "request does not contain an access token"})
		c.Abort()
		return nil
	}

	validateTokenErr, userData := auth.ValidateToken(authparts[1])
	if validateTokenErr != nil {
		c.JSON(401, gin.H{"error": "you must pass an Authorization header with a valid JWT bearer token"})
		c.Abort()
		return nil
	}

	return userData
}
