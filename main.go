package main

import (
	"github.com/dockrelix/dockrelix-backend/database"
	"github.com/dockrelix/dockrelix-backend/handlers"
	"github.com/dockrelix/dockrelix-backend/models"
	"github.com/dockrelix/dockrelix-backend/middleware"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	
	database.Connect()
	database.AutoMigrate()
	
	seedInitialUser()
	
	r := gin.Default()
	
	r.POST("/login", handlers.Login)
	r.GET("/protected", middleware.JWTAuth(), handlers.ProtectedEndpoint)
	
	log.Fatal(r.Run(":" + os.Getenv("PORT")))
}

func seedInitialUser() {
	var count int64
	database.DB.Model(&models.User{}).Count(&count)
	
	if count == 0 {
		user := models.User{
			Username: "admin",
			Password: handlers.HashPassword("admin"),
			Email:    "admin@example.com",
		}
		database.DB.Create(&user)
		log.Println("Initial user created")
	}
}
