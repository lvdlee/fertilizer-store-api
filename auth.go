package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lvdlee/fertilizer-store-data/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Login struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CreateJwt(user *models.User, secret string) (string, error) {
	validity, err := time.ParseDuration("1h")
	if err != nil {
		return "", err
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": user.ID,
		"exp":  time.Now().Add(validity).Unix(),
	})

	token, err := claims.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return token, nil
}

func ValidJwt(token, secret string) bool {
	parsed, err := jwt.Parse(token, func(_token *jwt.Token) (interface{}, error) {
		_, ok := _token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", _token.Header["alg"])
		}

		return []byte(secret), nil
	})

	if err != nil {
		return false
	}

	return parsed.Valid
}

func SetupAuthRoutes(r *gin.Engine, db *gorm.DB, secret string) {
	r.POST("/auth/login", func(c *gin.Context) {
		var body Login

		err := c.ShouldBindJSON(&body)
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

		user := &models.User{}

		db.Where("username = ?", body.Username).First(user)

		if user.Username != body.Username || !CheckPassword(body.Password, user.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid username or password",
			})
			return
		}

		jwt, err := CreateJwt(user, secret)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to create authentication token",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"access_token": jwt,
		})
	})
}
