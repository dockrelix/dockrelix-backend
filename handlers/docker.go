package handlers

import (
	"github.com/docker/docker/client"
	"github.com/dockrelix/dockrelix-backend/docker"

	"github.com/gin-gonic/gin"
)

func ListStacks(cli *client.Client, c *gin.Context) {
	result := docker.ListStacks(cli)
	c.JSON(200, result)
}
