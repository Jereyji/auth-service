package handlers_test

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Jereyji/auth-service/internal/auth/domain/entity"
	auth_errors "github.com/Jereyji/auth-service/internal/auth/domain/errors"
	"github.com/Jereyji/auth-service/internal/auth/presentation/handlers"
	handler_mock "github.com/Jereyji/auth-service/internal/auth/presentation/handlers/mocks"
	"github.com/Jereyji/auth-service/internal/pkg/configs"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	tokensCfg = &configs.TokensConfig{
		AccessTokenExpiresIn:  time.Minute * 15,
		RefreshTokenExpiresIn: time.Hour * 24 * 7,
	}

	logger = slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
)

func TestAuthHandler_Register_Success(t *testing.T) {
	mockService := new(handler_mock.MockIAuthService)
	mockKafka := new(handler_mock.MockIKafkaProducer)

	mockService.On("Register", mock.Anything, "testuser", "test@example.com", "password123").
		Return(nil)

	handler := handlers.NewAuthHandler(mockService, mockKafka, tokensCfg, logger)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	requestBody := `{"name":"testuser","email":"test@example.com","password":"password123"}`
	c.Request = httptest.NewRequest("POST", "/register", strings.NewReader(requestBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Register(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response handlers.RegisterResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "test@example.com", response.Email)

	mockService.AssertExpectations(t)
}

func TestAuthHandler_Register_UserExists(t *testing.T) {
	mockService := new(handler_mock.MockIAuthService)
	mockKafka := new(handler_mock.MockIKafkaProducer)

	mockService.On("Register", mock.Anything, "testuser", "existing@example.com", "password123").
		Return(auth_errors.ErrRowExist)

	handler := handlers.NewAuthHandler(mockService, mockKafka, tokensCfg, logger)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	requestBody := `{"name":"testuser","email":"existing@example.com","password":"password123"}`
	c.Request = httptest.NewRequest("POST", "/register", strings.NewReader(requestBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Register(c)

	assert.Equal(t, http.StatusConflict, w.Code)

	mockService.AssertExpectations(t)
}

func TestAuthHandler_Login_Success(t *testing.T) {
	mockService := new(handler_mock.MockIAuthService)
	mockKafka := new(handler_mock.MockIKafkaProducer)

	accessToken := entity.AccessToken{Token: "access_token_123"}
	refreshToken := entity.RefreshToken{Token: "refresh_token_456", ExpiresAt: time.Now().Add(time.Hour * 24)}

	mockService.On("Login", mock.Anything, "test@example.com", "password123").
		Return(accessToken, refreshToken, nil)

	mockKafka.On("SendMessage", "test@example.com", mock.Anything).
		Return(nil)

	handler := handlers.NewAuthHandler(mockService, mockKafka, tokensCfg, logger)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	requestBody := `{"email":"test@example.com","password":"password123"}`
	c.Request = httptest.NewRequest("POST", "/login", strings.NewReader(requestBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Login(c)

	assert.Equal(t, http.StatusOK, w.Code)

	cookies := w.Result().Cookies()
	var accessTokenCookie, refreshTokenCookie *http.Cookie

	for _, cookie := range cookies {
		if cookie.Name == "access_token" {
			accessTokenCookie = cookie
		}
		if cookie.Name == "refresh_token" {
			refreshTokenCookie = cookie
		}
	}

	assert.NotNil(t, accessTokenCookie)
	assert.NotNil(t, refreshTokenCookie)
	assert.Equal(t, "access_token_123", accessTokenCookie.Value)
	assert.Equal(t, "refresh_token_456", refreshTokenCookie.Value)

	mockService.AssertExpectations(t)
	mockKafka.AssertExpectations(t)
}

func TestAuthHandler_Login_InvalidCredentials(t *testing.T) {
	mockService := new(handler_mock.MockIAuthService)
	mockKafka := new(handler_mock.MockIKafkaProducer)

	mockService.On("Login", mock.Anything, "test@example.com", "wrongpassword").
		Return(entity.AccessToken{}, entity.RefreshToken{}, auth_errors.ErrInvalidEmailOrPassword)

	handler := handlers.NewAuthHandler(mockService, mockKafka, tokensCfg, logger)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	requestBody := `{"email":"test@example.com","password":"wrongpassword"}`
	c.Request = httptest.NewRequest("POST", "/login", strings.NewReader(requestBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Login(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	mockService.AssertExpectations(t)
}

func TestAuthHandler_RefreshTokens_Success(t *testing.T) {
	mockService := new(handler_mock.MockIAuthService)
	mockKafka := new(handler_mock.MockIKafkaProducer)

	oldRefreshToken := "old_refresh_token"
	newAccessToken := entity.AccessToken{Token: "new_access_token"}
	newRefreshToken := entity.RefreshToken{Token: "new_refresh_token", ExpiresAt: time.Now().Add(time.Hour * 24)}

	mockService.On("RefreshTokens", mock.Anything, oldRefreshToken).
		Return(newAccessToken, newRefreshToken, nil)

	handler := handlers.NewAuthHandler(mockService, mockKafka, tokensCfg, logger)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("POST", "/refresh", nil)
	c.Request.AddCookie(&http.Cookie{
		Name:  "refresh_token",
		Value: oldRefreshToken,
	})

	handler.RefreshTokens(c)

	assert.Equal(t, http.StatusOK, w.Code)

	cookies := w.Result().Cookies()
	var accessTokenCookie, refreshTokenCookie *http.Cookie

	for _, cookie := range cookies {
		if cookie.Name == "access_token" {
			accessTokenCookie = cookie
		}
		if cookie.Name == "refresh_token" {
			refreshTokenCookie = cookie
		}
	}

	assert.NotNil(t, accessTokenCookie)
	assert.NotNil(t, refreshTokenCookie)
	assert.Equal(t, "new_access_token", accessTokenCookie.Value)
	assert.Equal(t, "new_refresh_token", refreshTokenCookie.Value)

	mockService.AssertExpectations(t)
}

func TestAuthHandler_RefreshTokens_InvalidToken(t *testing.T) {
	mockService := new(handler_mock.MockIAuthService)
	mockKafka := new(handler_mock.MockIKafkaProducer)

	mockService.On("RefreshTokens", mock.Anything, "invalid_token").
		Return(entity.AccessToken{}, entity.RefreshToken{}, auth_errors.ErrNotFound)

	handler := handlers.NewAuthHandler(mockService, mockKafka, tokensCfg, logger)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("POST", "/refresh", nil)
	c.Request.AddCookie(&http.Cookie{
		Name:  "refresh_token",
		Value: "invalid_token",
	})

	handler.RefreshTokens(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	mockService.AssertExpectations(t)
}

// func TestAuthHandler_SendEventToKafka_Error(t *testing.T) {
// 	mockService := new(handler_mock.MockIAuthService)
// 	mockKafka := new(handler_mock.MockIKafkaProducer)

// 	accessToken := entity.AccessToken{Token: "access_token_123"}
// 	refreshToken := entity.RefreshToken{Token: "refresh_token_456", ExpiresAt: time.Now().Add(time.Hour * 24)}

// 	mockService.On("Login", mock.Anything, "test@example.com", "password123").
// 		Return(accessToken, refreshToken, nil)

// 	kafkaError := errors.New("kafka connection error")
// 	mockKafka.On("SendMessage", "test@example.com", mock.AnythingOfType("string")).
// 		Return(kafkaError)

// 	handler := handlers.NewAuthHandler(mockService, mockKafka, tokensCfg, logger)

// 	gin.SetMode(gin.TestMode)
// 	w := httptest.NewRecorder()
// 	c, _ := gin.CreateTestContext(w)

// 	originalStatus := c.Writer.Status()

// 	// Middleware для отслеживания изменений статуса
// 	c.Writer = &statusTrackingResponseWriter{ResponseWriter: c.Writer, t: t}

// 	requestBody := `{"email":"test@example.com","password":"password123"}`
// 	c.Request = httptest.NewRequest("POST", "/login", strings.NewReader(requestBody))
// 	c.Request.Header.Set("Content-Type", "application/json")

// 	// Вызываем обработчик
// 	handler.Login(c)

// 	// Проверяем статус сразу после вызова Login
// 	afterLoginStatus := c.Writer.Status()

// 	t.Logf("Original status: %d, After Login status: %d, Response code: %d",
// 		originalStatus, afterLoginStatus, w.Code)

// 	// Проверяем финальный статус в ResponseRecorder
// 	assert.Equal(t, http.StatusInternalServerError, w.Code,
// 		"Ожидался код 500, но получен %d", w.Code)

// 	mockService.AssertExpectations(t)
// 	mockKafka.AssertExpectations(t)
// }

// // Вспомогательный тип для отслеживания изменений статуса
// type statusTrackingResponseWriter struct {
// 	gin.ResponseWriter
// 	t *testing.T
// }

// func (w *statusTrackingResponseWriter) WriteHeader(code int) {
// 	w.t.Logf("Setting status code to: %d", code)
// 	w.ResponseWriter.WriteHeader(code)
// }
