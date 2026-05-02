package api

import (
	"net/http"

	"the-unified-document-viewer/internal/auth"
	"the-unified-document-viewer/internal/models"
	"the-unified-document-viewer/internal/repository"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	JWTManager *auth.JWTManager
	UserRepo   *repository.UserRepository
}

func NewAuthHandler(jwtManager *auth.JWTManager, userRepo *repository.UserRepository) *AuthHandler {
	return &AuthHandler{
		JWTManager: jwtManager,
		UserRepo:   userRepo,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	user, err := h.UserRepo.FindByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !auth.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := h.JWTManager.GenerateToken(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	refreshToken, err := h.JWTManager.GenerateRefreshToken(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(h.JWTManager.TokenExpiry().Seconds()),
	})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req models.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	claims, err := h.JWTManager.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// Generate new access token
	token, err := h.JWTManager.GenerateToken(claims.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{
		Token:      token,
		ExpiresIn:  int64(h.JWTManager.TokenExpiry().Seconds()),
	})
}
