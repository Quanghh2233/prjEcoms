package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/qhh/prjEcom/pkg/api/middlewares"
	"github.com/qhh/prjEcom/pkg/db/store"
	"github.com/qhh/prjEcom/pkg/models"
)

type ProductHandler struct {
	store *store.Store
}

func NewProductHandler(store *store.Store) *ProductHandler {
	return &ProductHandler{
		store: store,
	}
}

type createProductRequest struct {
	ShopID      string   `json:"shop_id" binding:"required"`
	Name        string   `json:"name" binding:"required,min=3,max=100"`
	Description string   `json:"description" binding:"required"`
	Price       float64  `json:"price" binding:"required,gt=0"`
	Stock       int32    `json:"stock" binding:"required,min=0"`
	Category    string   `json:"category" binding:"required"`
	ImageURLs   []string `json:"image_urls"`
}

// CreateProduct creates a new product
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req createProductRequest
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

	// Parse shop ID
	shopID, err := uuid.Parse(req.ShopID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shop ID"})
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
		c.JSON(http.StatusForbidden, gin.H{"error": "you don't have permission to create products in this shop"})
		return
	}

	// Create product
	product := &models.Product{
		ShopID:      shopID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Category:    req.Category,
		ImageURLs:   req.ImageURLs,
	}

	err = h.store.CreateProduct(c, product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create product"})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// GetProduct returns a single product by ID
func (h *ProductHandler) GetProduct(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	product, err := h.store.GetProductByID(c, productID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// ListProducts returns a paginated list of all products
func (h *ProductHandler) ListProducts(c *gin.Context) {
	var req struct {
		Limit  int `form:"limit" binding:"required,min=1,max=100"`
		Offset int `form:"offset" binding:"required,min=0"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	products, err := h.store.ListProducts(c, req.Limit, req.Offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list products"})
		return
	}

	c.JSON(http.StatusOK, products)
}

// SearchProducts searches for products by name or description
func (h *ProductHandler) SearchProducts(c *gin.Context) {
	var req struct {
		Query  string `form:"q" binding:"required,min=1"`
		Limit  int    `form:"limit" binding:"required,min=1,max=100"`
		Offset int    `form:"offset" binding:"required,min=0"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	products, err := h.store.SearchProducts(c, req.Query, req.Limit, req.Offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to search products"})
		return
	}

	c.JSON(http.StatusOK, products)
}

// ListProductsByShop returns products from a specific shop
func (h *ProductHandler) ListProductsByShop(c *gin.Context) {
	shopID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shop ID"})
		return
	}

	var req struct {
		Limit  int `form:"limit" binding:"required,min=1,max=100"`
		Offset int `form:"offset" binding:"required,min=0"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	products, err := h.store.GetProductsByShopID(c, shopID, req.Limit, req.Offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get shop products"})
		return
	}

	c.JSON(http.StatusOK, products)
}

// FilterProductsByCategory returns products filtered by category
func (h *ProductHandler) FilterProductsByCategory(c *gin.Context) {
	category := c.Query("category")
	if category == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category is required"})
		return
	}

	var req struct {
		Limit  int `form:"limit" binding:"required,min=1,max=100"`
		Offset int `form:"offset" binding:"required,min=0"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	products, err := h.store.FilterProductsByCategory(c, category, req.Limit, req.Offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to filter products"})
		return
	}

	c.JSON(http.StatusOK, products)
}

// FilterProductsByPrice returns products filtered by price range
func (h *ProductHandler) FilterProductsByPrice(c *gin.Context) {
	minPriceStr := c.Query("min_price")
	maxPriceStr := c.Query("max_price")

	if minPriceStr == "" || maxPriceStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "min_price and max_price are required"})
		return
	}

	minPrice, err := strconv.ParseFloat(minPriceStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid min_price"})
		return
	}

	maxPrice, err := strconv.ParseFloat(maxPriceStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid max_price"})
		return
	}

	if minPrice < 0 || maxPrice < minPrice {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid price range"})
		return
	}

	var req struct {
		Limit  int `form:"limit" binding:"required,min=1,max=100"`
		Offset int `form:"offset" binding:"required,min=0"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	products, err := h.store.FilterProductsByPrice(c, minPrice, maxPrice, req.Limit, req.Offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to filter products by price"})
		return
	}

	c.JSON(http.StatusOK, products)
}

// UpdateProduct updates a product's details
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	var req createProductRequest
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

	// Get product to check ownership
	product, err := h.store.GetProductByID(c, productID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	// Get shop to check ownership
	shop, err := h.store.GetShopByID(c, product.ShopID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get shop"})
		return
	}

	// Check if user owns the shop or is admin
	if shop.UserID != payload.UserID && payload.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "you don't have permission to update this product"})
		return
	}

	// Update product
	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	product.Stock = req.Stock
	product.Category = req.Category
	product.ImageURLs = req.ImageURLs

	err = h.store.UpdateProduct(c, product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update product"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// DeleteProduct deletes a product
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	// Get user from auth payload
	payload, err := middlewares.GetAuthPayload(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Get product to check ownership
	product, err := h.store.GetProductByID(c, productID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	// Get shop to check ownership
	shop, err := h.store.GetShopByID(c, product.ShopID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get shop"})
		return
	}

	// Check if user owns the shop or is admin
	if shop.UserID != payload.UserID && payload.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "you don't have permission to delete this product"})
		return
	}

	// Delete product
	err = h.store.DeleteProduct(c, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product deleted successfully"})
}
