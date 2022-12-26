package auth

import (
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var jwtKey = []byte("supersecretkey")

type JWTClaim struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
	jwt.StandardClaims
}

func GenerateJWT(name string, email string, age int) (tokenString string, err error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &JWTClaim{
		Email:    email,
		Username: name,
		Age:      age,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(jwtKey)

	return tokenString, err
}

func ValidateToken(token string) (error, *JWTClaim) {
	parsedToken, err := jwt.ParseWithClaims(token,
		&JWTClaim{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		})

	if err != nil {
		return err, nil
	}

	claims, ok := parsedToken.Claims.(*JWTClaim)

	if !ok {
		err = errors.New("couldn't parse claims")
		return err, nil
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("token expired")
		return err, nil
	}

	return nil, claims
}

func Auth() gin.HandlerFunc {
	return func(context *gin.Context) {
		authHeader := context.GetHeader("Authorization")
		if authHeader == "" {
			context.JSON(401, gin.H{"error": "request does not contain an authorization header"})
			context.Abort()
			return
		}

		authparts := strings.Split(authHeader, " ")

		if len(authparts) < 2 {
			context.JSON(401, gin.H{"error": "request does not contain an access token"})
			context.Abort()
			return
		}

		validateTokenErr, _ := ValidateToken(authparts[1])
		if validateTokenErr != nil {
			context.JSON(401, gin.H{"error": "you must pass an Authorization header with a valid JWT bearer token"})
			context.Abort()
			return
		}
	}
}
