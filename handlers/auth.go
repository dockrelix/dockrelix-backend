package handlers

import (
	"strings"
	"time"

	"github.com/dockrelix/dockrelix-backend/database"
	"github.com/dockrelix/dockrelix-backend/models"
	"github.com/dockrelix/dockrelix-backend/utils"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes)
}

func Login(c *gin.Context) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", utils.SanitizeInput(credentials.Email)).First(&user).Error; err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	c.JSON(200, gin.H{"token": tokenString})
}

func IsSetup(c *gin.Context) {
	var count int64
	database.DB.Model(&models.User{}).Count(&count)

	c.JSON(200, gin.H{"is_setup": count == 0})
}

func Setup(c *gin.Context) {
	var information struct {
		Username     string `json:"username"`
		Password     string `json:"password"`
		Email        string `json:"email"`
		Organization string `json:"organization"`
	}

	if err := c.ShouldBindJSON(&information); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	if information.Username == "" || information.Password == "" || information.Email == "" || information.Organization == "" {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	information.Username = utils.SanitizeInput(information.Username)
	information.Email = utils.SanitizeInput(information.Email)
	information.Organization = strings.ToLower(utils.SanitizeInput(information.Organization))

	if !utils.IsAlphaNumeric(information.Organization) {
		c.JSON(400, gin.H{"error": "Invalid organization name"})
		return
	}

	if !utils.IsEmailValid(information.Email) {
		c.JSON(400, gin.H{"error": "Invalid email"})
		return
	}

	if len(information.Password) < 8 {
		c.JSON(400, gin.H{"error": "Password must be at least 8 characters long"})
		return
	}

	hashedPassword := HashPassword(information.Password)

	user := models.User{
		Username:     information.Username,
		Password:     hashedPassword,
		Email:        information.Email,
		Organization: information.Organization,
	}

	database.DB.Create(&user)

	c.JSON(200, gin.H{"message": "Setup complete"})
}
