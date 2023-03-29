package main

import (
	"context"

	"github.com/gin-gonic/gin"
	database "github.com/lvdlee/fertilizer-store-data"
)

func main() {
	ctx := context.Background()

	settings := ReadSettings()

	println("[DEBUG] Setting up database")
	db := database.Setup(ctx, settings.Database.Location)
	println("[DEBUG] Database connection established")

	r := gin.Default()

	SetupAuthRoutes(r, db, settings.Auth.Jwt.Secret)
	SetupUsersRoutes(r, db)

	r.Run("0.0.0.0:3000")
}
