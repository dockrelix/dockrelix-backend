package handlers

import (
	"github.com/docker/docker/client"
	"github.com/dockrelix/dockrelix-backend/docker"
	"github.com/dockrelix/dockrelix-backend/models"

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

func CreateStackDraft(cli *client.Client, c *gin.Context) {
	var stackDraft models.StackDraft
	if err := c.ShouldBindJSON(&stackDraft); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if stackDraft.Name == "" {
		c.JSON(400, gin.H{"error": "Name is required"})
		return
	}

	if stackDraft.Data == "" {
		c.JSON(400, gin.H{"error": "Data is required"})
		return
	}

	err := docker.SaveDraft(stackDraft)
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: stack_drafts.name" {
			c.JSON(400, gin.H{"error": "Draft with this name already exists"})
			return
		}
		c.JSON(500, gin.H{"error": "Stack could not be drafted: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Stack draft created successfully"})
}

func GetStackDrafts(c *gin.Context) {
	result := docker.GetDrafts()
	c.JSON(200, result)
}
