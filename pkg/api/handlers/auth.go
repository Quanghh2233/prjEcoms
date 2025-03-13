package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
	"github.com/qhh/prjEcom/pkg/db/store"
	"github.com/qhh/prjEcom/pkg/models"
	"github.com/qhh/prjEcom/pkg/utils"
)

type AuthHandler struct {
	store    *store.Store
	jwtMaker *utils.JWTMaker
}

func NewAuthHandler(store *store.Store, jwtMaker *utils.JWTMaker) *AuthHandler {
	return &AuthHandler{
		store:    store,
		jwtMaker: jwtMaker,
	}
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Token    string `json:"token"`
}

type registerRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if username already exists
	_, err := h.store.GetUserByUsername(c, req.Username)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username already exists"})
		return
	} else if !errors.Is(err, pg.ErrNoRows) && err.Error() != "user not found" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check username"})
		return
	}

	// Check if email already exists
	_, err = h.store.GetUserByEmail(c, req.Email)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email already exists"})
		return
	} else if !errors.Is(err, pg.ErrNoRows) && err.Error() != "user not found" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check email"})
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	// Create user
	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Role:         models.RoleBuyer,
	}

	err = h.store.CreateUser(c, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
	})
}

// Login handles user authentication
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user by username
	user, err := h.store.GetUserByUsername(c, req.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Check if user is banned
	if user.IsBanned {
		c.JSON(http.StatusForbidden, gin.H{"error": "your account has been banned"})
		return
	}

	// Verify password
	err = utils.CheckPassword(req.Password, user.PasswordHash)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Generate JWT token
	token, err := h.jwtMaker.CreateToken(user.ID, user.Username, string(user.Role), 24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create token"})
		return
	}

	c.JSON(http.StatusOK, loginResponse{
		ID:       user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
		Role:     string(user.Role),
		Token:    token,
	})
}
