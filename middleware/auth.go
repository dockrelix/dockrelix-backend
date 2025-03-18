package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/dockrelix/dockrelix-backend/database"
	"github.com/dockrelix/dockrelix-backend/models"
	"os"
	"strings"
	"fmt"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Authorization header required"})
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		fmt.Println(tokenString)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			var user models.User
			if err := database.DB.First(&user, claims["sub"]).Error; err != nil {
				c.AbortWithStatusJSON(401, gin.H{"error": "User not found"})
				return
			}
			c.Set("user", user)
			c.Next()
		} else {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token", "details": err.Error()})
		}
	}
}
