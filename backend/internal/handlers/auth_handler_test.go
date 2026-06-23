package handlers_test

import (
	"booktracker/backend/internal/handlers"
	"booktracker/backend/internal/repositories"
	"booktracker/backend/internal/services"
	"booktracker/backend/internal/testutil"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupAuthRouter(t *testing.T) (*gin.Engine, func()) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	db := testutil.NewTestDB(t)
	testutil.CleanDB(t, db)

	authRepo := repositories.NewAuthRepository(db)
	authService := services.NewAuthService(authRepo)
	authHandler := handlers.NewAuthHandler(authService)

	r := gin.New()
	r.POST("/api/auth/register", authHandler.Register)
	r.POST("/api/auth/login", authHandler.Login)
	r.POST("/api/auth/refresh", authHandler.Refresh)
	r.POST("/api/auth/logout", authHandler.Logout)

	return r, func() { db.Close() }
}

func TestRegister_Success(t *testing.T) {
	r, cleanup := setupAuthRouter(t)
	defer cleanup()

	body := map[string]string{
		"name":     "Test User",
		"email":    "test@example.com",
		"password": "password123",
		"timezone": "UTC",
	}
	bodyJSON, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got: %d body: %s", w.Code, w.Body.String())
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["token"] == nil {
		t.Error("expected token in response")
	}
	user, ok := response["user"].(map[string]interface{})
	if !ok {
		t.Fatal("expected user in response")
	}
	if user["email"] != "test@example.com" {
		t.Errorf("expected email test@example.com, got: %v", user["email"])
	}
}

func TestRegister_DuplicateEmail(t *testing.T) {
	r, cleanup := setupAuthRouter(t)
	defer cleanup()

	body := map[string]string{
		"name":     "Test User",
		"email":    "test@example.com",
		"password": "password123",
	}
	bodyJSON, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("first register failed: %d", w.Code)
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("expected status 409, got: %d", w.Code)
	}
}

func TestRegister_InvalidEmail(t *testing.T) {
	r, cleanup := setupAuthRouter(t)
	defer cleanup()

	body := map[string]string{
		"name":     "Test User",
		"email":    "notanemail",
		"password": "password123",
	}
	bodyJSON, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got: %d", w.Code)
	}
}

func TestRegister_ShortPassword(t *testing.T) {
	r, cleanup := setupAuthRouter(t)
	defer cleanup()

	body := map[string]string{
		"name":     "Test User",
		"email":    "test@example.com",
		"password": "short",
	}
	bodyJSON, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got: %d", w.Code)
	}
}

func TestLogin_Success(t *testing.T) {
	r, cleanup := setupAuthRouter(t)
	defer cleanup()

	regBody := map[string]string{
		"name":     "Test User",
		"email":    "test@example.com",
		"password": "password123",
	}
	regJSON, _ := json.Marshal(regBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(regJSON))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("register failed: %d %s", w.Code, w.Body.String())
	}

	loginBody := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	loginJSON, _ := json.Marshal(loginBody)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(loginJSON))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got: %d body: %s", w.Code, w.Body.String())
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["token"] == nil {
		t.Error("expected token in response")
	}

	cookies := w.Result().Cookies()
	found := false
	for _, c := range cookies {
		if c.Name == "refresh_token" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected refresh_token cookie to be set")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	r, cleanup := setupAuthRouter(t)
	defer cleanup()

	regBody := map[string]string{
		"name":     "Test User",
		"email":    "test@example.com",
		"password": "password123",
	}
	regJSON, _ := json.Marshal(regBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(regJSON))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	loginBody := map[string]string{
		"email":    "test@example.com",
		"password": "wrongpassword",
	}
	loginJSON, _ := json.Marshal(loginBody)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(loginJSON))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got: %d", w.Code)
	}
}

func TestLogin_NonExistentEmail(t *testing.T) {
	r, cleanup := setupAuthRouter(t)
	defer cleanup()

	loginBody := map[string]string{
		"email":    "nobody@example.com",
		"password": "password123",
	}
	loginJSON, _ := json.Marshal(loginBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(loginJSON))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got: %d", w.Code)
	}
}

func TestLogout_Success(t *testing.T) {
	r, cleanup := setupAuthRouter(t)
	defer cleanup()

	regBody := map[string]string{
		"name":     "Test User",
		"email":    "test@example.com",
		"password": "password123",
	}
	regJSON, _ := json.Marshal(regBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(regJSON))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	var regResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &regResponse)
	token, ok := regResponse["token"].(string)
	if !ok || token == "" {
		t.Fatalf("expected token in register response, got: %s", w.Body.String())
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/auth/logout", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got: %d", w.Code)
	}

	cookies := w.Result().Cookies()
	for _, c := range cookies {
		if c.Name == "refresh_token" && c.MaxAge > 0 {
			t.Error("expected refresh_token cookie to be cleared")
		}
	}
}