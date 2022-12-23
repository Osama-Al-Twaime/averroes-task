package handlers

import (
	"averroes-task/models"
	"fmt"

	"github.com/gin-gonic/gin"
)

func RegisterNewUser(c *gin.Context) {
	var userREquestBody models.UserRegisterRequest

	if err := c.BindJSON(&userREquestBody); err != nil {
		fmt.Println("Error with the request body")
	}

	fmt.Println(userREquestBody)
}
