package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"github.com/qhh/prjEcom/pkg/api/middlewares"
	"github.com/qhh/prjEcom/pkg/db/store"
	"github.com/qhh/prjEcom/pkg/models"
)

type OrderHandler struct {
	store *store.Store
}

func NewOrderHandler(store *store.Store) *OrderHandler {
	return &OrderHandler{
		store: store,
	}
}

type orderItemRequest struct {
	ProductID string  `json:"product_id" binding:"required"`
	Quantity  int32   `json:"quantity" binding:"required,min=1"`
	Price     float64 `json:"price" binding:"required,gt=0"`
}

type createOrderRequest struct {
	ShippingAddress string             `json:"shipping_address" binding:"required"`
	Items           []orderItemRequest `json:"items" binding:"required,min=1"`
}

// CreateOrder creates a new order
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req createOrderRequest
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

	// Calculate total amount
	var totalAmount float64
	for _, item := range req.Items {
		totalAmount += item.Price * float64(item.Quantity)
	}

	// Prepare order
	order := &models.Order{
		UserID:          payload.UserID,
		TotalAmount:     totalAmount,
		ShippingAddress: req.ShippingAddress,
		Status:          models.StatusPending,
	}

	// Start transaction
	var orderItems []*models.OrderItem
	err = h.store.RunInTransaction(c, func(tx *pg.Tx) error {
		// Create order
		if err := tx.ModelContext(c, order).Insert(); err != nil {
			return err
		}

		// Process order items
		for _, item := range req.Items {
			productID, err := uuid.Parse(item.ProductID)
			if err != nil {
				return err
			}

			// Check if product exists and has enough stock
			product := &models.Product{ID: productID}
			if err := tx.ModelContext(c, product).WherePK().Select(); err != nil {
				return err
			}

			if product.Stock < item.Quantity {
				return &stockError{
					ProductName: product.Name,
					Stock:       product.Stock,
					Requested:   item.Quantity,
				}
			}

			// Add order item
			orderItem := &models.OrderItem{
				OrderID:         order.ID,
				ProductID:       productID,
				Quantity:        item.Quantity,
				PriceAtPurchase: item.Price,
			}
			if err := tx.ModelContext(c, orderItem).Insert(); err != nil {
				return err
			}
			orderItems = append(orderItems, orderItem)

			// Update product stock
			product.Stock -= item.Quantity
			if _, err := tx.ModelContext(c, product).WherePK().Update(); err != nil {
				return err
			}
		}

		return nil
	})

	// Handle errors
	if err != nil {
		if stockErr, ok := err.(*stockError); ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": stockErr.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to create order: " + err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"order_id":     order.ID,
		"total_amount": order.TotalAmount,
		"status":       order.Status,
		"created_at":   order.CreatedAt,
	})
}

// GetOrder returns a specific order
func (h *OrderHandler) GetOrder(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}

	// Get user from auth payload
	payload, err := middlewares.GetAuthPayload(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Get order
	order, err := h.store.GetOrderByID(c, orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	// Check if user owns the order or is admin
	if order.UserID != payload.UserID && payload.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "you don't have permission to view this order"})
		return
	}

	// Get order items
	items, err := h.store.GetOrderItems(c, orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get order items"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"order": order,
		"items": items,
	})
}

// GetUserOrders returns orders for the current user
func (h *OrderHandler) GetUserOrders(c *gin.Context) {
	// Get user from auth payload
	payload, err := middlewares.GetAuthPayload(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
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

	// Get user orders
	orders, err := h.store.GetOrdersByUserID(c, payload.UserID, req.Limit, req.Offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// UpdateOrderStatus updates an order's status (admin only)
func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required,oneof=pending paid shipped delivered canceled"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update order status
	order, err := h.store.UpdateOrderStatus(c, orderID, models.OrderStatus(req.Status))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update order status"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// stockError is a custom error type for insufficient stock
type stockError struct {
	ProductName string
	Stock       int32
	Requested   int32
}

func (e *stockError) Error() string {
	return fmt.Sprintf("insufficient stock for product: %s (available: %d, requested: %d)",
		e.ProductName, e.Stock, e.Requested)
}
