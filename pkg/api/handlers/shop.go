package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/qhh/prjEcom/pkg/api/middlewares"
	"github.com/qhh/prjEcom/pkg/db/store"
	"github.com/qhh/prjEcom/pkg/models"
)

type ShopHandler struct {
	store *store.Store
}

func NewShopHandler(store *store.Store) *ShopHandler {
	return &ShopHandler{
		store: store,
	}
}

type createShopRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=100"`
	Description string `json:"description" binding:"required"`
	LogoURL     string `json:"logo_url"`
}

// CreateShop creates a new shop for the authenticated user
func (h *ShopHandler) CreateShop(c *gin.Context) {
	var req createShopRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user from auth payload
	payload, err := middlewares.GetAuthPayload(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Create shop
	shop := &models.Shop{
		UserID:      payload.UserID,
		Name:        req.Name,
		Description: req.Description,
		LogoURL:     req.LogoURL,
	}

	err = h.store.CreateShop(c, shop)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create shop"})
		return
	}

	// Update user's role to seller
	_, err = h.store.GetUserByID(c, payload.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user"})
		return
	}

	_, err = h.store.UpdateUserRole(c, payload.UserID, models.RoleSeller)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user role"})
		return
	}

	c.JSON(http.StatusCreated, shop)
}

// GetShop returns details of a specific shop
func (h *ShopHandler) GetShop(c *gin.Context) {
	shopID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shop ID"})
		return
	}

	shop, err := h.store.GetShopByID(c, shopID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "shop not found"})
		return
	}

	c.JSON(http.StatusOK, shop)
}

// GetUserShops returns all shops owned by the authenticated user
func (h *ShopHandler) GetUserShops(c *gin.Context) {
	// Get user from auth payload
	payload, err := middlewares.GetAuthPayload(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	shops, err := h.store.GetShopsByUserID(c, payload.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get shops"})
		return
	}

	c.JSON(http.StatusOK, shops)
}

// ListShops returns a paginated list of all shops
func (h *ShopHandler) ListShops(c *gin.Context) {
	var req struct {
		Limit  int `form:"limit" binding:"required,min=1,max=100"`
		Offset int `form:"offset" binding:"required,min=0"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shops, err := h.store.ListShops(c, req.Limit, req.Offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list shops"})
		return
	}

	c.JSON(http.StatusOK, shops)
}

// SearchShops searches for shops by name or description
func (h *ShopHandler) SearchShops(c *gin.Context) {
	var req struct {
		Query  string `form:"q" binding:"required,min=1"`
		Limit  int    `form:"limit" binding:"required,min=1,max=100"`
		Offset int    `form:"offset" binding:"required,min=0"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shops, err := h.store.SearchShops(c, req.Query, req.Limit, req.Offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to search shops"})
		return
	}

	c.JSON(http.StatusOK, shops)
}

// UpdateShop updates a shop's details
func (h *ShopHandler) UpdateShop(c *gin.Context) {
	shopID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shop ID"})
		return
	}

	var req createShopRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user from auth payload
	payload, err := middlewares.GetAuthPayload(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Check if shop exists and belongs to user
	shop, err := h.store.GetShopByID(c, shopID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "shop not found"})
		return
	}

	// Check if user owns the shop or is admin
	if shop.UserID != payload.UserID && payload.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "you don't have permission to update this shop"})
		return
	}

	// Update shop
	shop.Name = req.Name
	shop.Description = req.Description
	shop.LogoURL = req.LogoURL

	err = h.store.UpdateShop(c, shop)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update shop"})
		return
	}

	c.JSON(http.StatusOK, shop)
}
