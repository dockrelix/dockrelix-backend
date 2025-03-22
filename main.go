package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dockrelix/dockrelix-backend/database"
	"github.com/dockrelix/dockrelix-backend/docker"
	"github.com/dockrelix/dockrelix-backend/handlers"
	"github.com/dockrelix/dockrelix-backend/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	fmt.Println("trying to connect to database")
	database.Connect()
	database.AutoMigrate()

	cli, err := docker.NewDockerClient()
	if err != nil {
		log.Fatalf("Failed to create Docker client: %v", err)
	}
	defer cli.Close()

	r := gin.Default()

	auth := r.Group("/auth")
	{
		auth.POST("/login", handlers.Login)

		// Setup routes
		auth.GET("/is-setup", handlers.IsSetup)
		auth.POST("/setup", handlers.Setup)
	}

	docker := r.Group("/docker")
	docker.Use(middleware.JWTAuth())
	{
		docker.GET("/stacks", func(c *gin.Context) {
			handlers.ListStacks(cli, c)
		})

		docker.GET("/stacks/:name", func(c *gin.Context) {
			handlers.ParseStackConfig(cli, c)
		})

		docker.POST("/stacks/draft", func(c *gin.Context) {
			handlers.CreateStackDraft(cli, c)
		})

		docker.GET("/stacks/drafts", func(c *gin.Context) {
			handlers.GetStackDrafts(c)
		})
	}

	log.Fatal(r.Run(":" + os.Getenv("PORT")))
}
