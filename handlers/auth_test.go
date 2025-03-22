package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"os"

	"github.com/dockrelix/dockrelix-backend/database"
	"github.com/dockrelix/dockrelix-backend/handlers"
	"github.com/dockrelix/dockrelix-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.POST("/login", handlers.Login)
	router.GET("/is-setup", handlers.IsSetup)
	router.POST("/setup", handlers.Setup)

	return router
}

func TestLogin(t *testing.T) {
	database.InitDBForTesting()

	router := setupRouter()

	user := models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: handlers.HashPassword("testpassword"),
	}

	if err := database.DB.Create(&user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	loginPayload := map[string]string{
		"email":    "test@example.com",
		"password": "testpassword",
	}
	payloadBytes, _ := json.Marshal(loginPayload)
	req, _ := http.NewRequest("POST", "/login", bytes.NewReader(payloadBytes))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %v", w.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	tokenString, ok := response["token"]
	if !ok {
		t.Fatal("missing token in response")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		t.Fatalf("error parsing token: %v", err)
	}

	if token == nil {
		t.Fatal("parsed token is nil")
	}
}

func TestLoginInvalidCredentials(t *testing.T) {
	router := setupRouter()

	database.InitDBForTesting()

	loginPayload := map[string]string{
		"email":    "wrong@example.com",
		"password": "wrongpassword",
	}
	payloadBytes, _ := json.Marshal(loginPayload)
	req, _ := http.NewRequest("POST", "/login", bytes.NewReader(payloadBytes))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %v", w.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if response["error"] != "Invalid credentials" {
		t.Errorf("expected 'Invalid credentials' error, got %v", response["error"])
	}
}

func TestIsSetup(t *testing.T) {
	router := setupRouter()

	database.InitDBForTesting()

	var count int64
	database.DB.Model(&models.User{}).Count(&count)
	if count != 0 {
		t.Fatalf("expected 0 users, found %v", count)
	}

	req, _ := http.NewRequest("GET", "/is-setup", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var response map[string]bool
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if !response["is_setup"] {
		t.Errorf("expected is_setup to be true, got false")
	}
}

func TestSetup(t *testing.T) {
	router := setupRouter()

	database.InitDBForTesting()

	var count int64
	database.DB.Model(&models.User{}).Count(&count)
	if count != 0 {
		t.Fatalf("expected 0 users, found %v", count)
	}

	setupPayload := map[string]string{
		"username":     "admin",
		"password":     "password123",
		"email":        "admin@example.com",
		"organization": "dockrelix",
	}

	payloadBytes, _ := json.Marshal(setupPayload)
	req, _ := http.NewRequest("POST", "/setup", bytes.NewReader(payloadBytes))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %v", w.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if response["message"] != "Setup complete" {
		t.Errorf("expected 'Setup complete', got %v", response["message"])
	}

	var newUser models.User
	if err := database.DB.Where("email = ?", setupPayload["email"]).First(&newUser).Error; err != nil {
		t.Fatalf("user not found in the database: %v", err)
	}

	if newUser.Username != setupPayload["username"] {
		t.Errorf("expected username %v, got %v", setupPayload["username"], newUser.Username)
	}
}

func TestSetupAlreadyComplete(t *testing.T) {
	database.InitDBForTesting()
	user := models.User{
		Username: "admin",
		Email:    "admin@example.com",
		Password: handlers.HashPassword("password123"),
	}
	database.DB.Create(&user)

	router := setupRouter()
	setupPayload := map[string]string{
		"username":     "admin",
		"password":     "password123",
		"email":        "admin@example.com",
		"organization": "dockrelix",
	}

	payloadBytes, _ := json.Marshal(setupPayload)
	req, _ := http.NewRequest("POST", "/setup", bytes.NewReader(payloadBytes))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %v", w.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if response["error"] != "Setup already complete" {
		t.Errorf("expected 'Setup already complete' error, got %v", response["error"])
	}
}
