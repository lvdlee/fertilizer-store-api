package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lvdlee/fertilizer-store-data/models"
	"gorm.io/gorm"
)

type UserCreate struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func SetupUsersRoutes(r *gin.Engine, db *gorm.DB) {
	r.GET("/users", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/users", func(c *gin.Context) {
		var request UserCreate

		err := c.ShouldBindJSON(&request)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid request body",
				"example": gin.H{
					"username": "johndoe",
					"password": "secret123",
				},
			})
			return
		}

		if len([]byte(request.Password)) > 72 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Password too long",
			})
			return
		}

		hashedPassword, err := HashPassword(request.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "An error occurred while hashing your password",
			})
			return
		}

		user := &models.User{}

		db.Where("username = ?", request.Username).First(user)

		if user.Username == request.Username {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "Username already exists",
			})
			return
		}

		user = &models.User{
			Username: request.Username,
			Password: hashedPassword,
		}

		db.Create(user)

		c.JSON(http.StatusOK, gin.H{
			"id":       user.ID,
			"username": request.Username,
		})
	})
}
