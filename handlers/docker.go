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

func ParseStackConfig(cli *client.Client, c *gin.Context) {
	stackName := c.Param("name")
	result, err := docker.ParseStackConfig(cli, stackName)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, result)
}
