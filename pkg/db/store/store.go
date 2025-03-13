package store

import (
	"context"
	"errors"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"github.com/qhh/prjEcom/pkg/models"
)

// Store provides all functions to execute db operations
type Store struct {
	db *pg.DB
}

// NewStore creates a new Store
func NewStore(db *pg.DB) *Store {
	return &Store{
		db: db,
	}
}

// User operations
func (s *Store) CreateUser(ctx context.Context, user *models.User) error {
	_, err := s.db.ModelContext(ctx, user).Insert()
	return err
}

func (s *Store) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user := &models.User{ID: id}
	err := s.db.ModelContext(ctx, user).WherePK().Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}

func (s *Store) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	err := s.db.ModelContext(ctx, user).Where("email = ?", email).Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}

func (s *Store) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	user := &models.User{}
	err := s.db.ModelContext(ctx, user).Where("username = ?", username).Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}

func (s *Store) ListUsers(ctx context.Context, limit, offset int) ([]*models.User, error) {
	var users []*models.User
	err := s.db.ModelContext(ctx, &users).
		Order("created_at ASC").
		Limit(limit).
		Offset(offset).
		Select()
	return users, err
}

func (s *Store) UpdateUserRole(ctx context.Context, id uuid.UUID, role models.UserRole) (*models.User, error) {
	user := &models.User{ID: id}
	err := s.db.ModelContext(ctx, user).WherePK().Select()
	if err != nil {
		return nil, err
	}

	user.Role = role
	_, err = s.db.ModelContext(ctx, user).WherePK().Update()
	return user, err
}

func (s *Store) BanUser(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user := &models.User{ID: id}
	err := s.db.ModelContext(ctx, user).WherePK().Select()
	if err != nil {
		return nil, err
	}

	user.IsBanned = true
	_, err = s.db.ModelContext(ctx, user).WherePK().Update()
	return user, err
}

func (s *Store) UnbanUser(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user := &models.User{ID: id}
	err := s.db.ModelContext(ctx, user).WherePK().Select()
	if err != nil {
		return nil, err
	}

	user.IsBanned = false
	_, err = s.db.ModelContext(ctx, user).WherePK().Update()
	return user, err
}

// Shop operations
func (s *Store) CreateShop(ctx context.Context, shop *models.Shop) error {
	_, err := s.db.ModelContext(ctx, shop).Insert()
	return err
}

func (s *Store) GetShopByID(ctx context.Context, id uuid.UUID) (*models.Shop, error) {
	shop := &models.Shop{ID: id}
	err := s.db.ModelContext(ctx, shop).WherePK().Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, errors.New("shop not found")
		}
		return nil, err
	}
	return shop, nil
}

func (s *Store) GetShopsByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Shop, error) {
	var shops []*models.Shop
	err := s.db.ModelContext(ctx, &shops).
		Where("user_id = ?", userID).
		Select()
	return shops, err
}

func (s *Store) ListShops(ctx context.Context, limit, offset int) ([]*models.Shop, error) {
	var shops []*models.Shop
	err := s.db.ModelContext(ctx, &shops).
		Order("created_at ASC").
		Limit(limit).
		Offset(offset).
		Select()
	return shops, err
}

func (s *Store) SearchShops(ctx context.Context, query string, limit, offset int) ([]*models.Shop, error) {
	var shops []*models.Shop
	err := s.db.ModelContext(ctx, &shops).
		Where("name ILIKE ? OR description ILIKE ?", "%"+query+"%", "%"+query+"%").
		Order("created_at ASC").
		Limit(limit).
		Offset(offset).
		Select()
	return shops, err
}

func (s *Store) UpdateShop(ctx context.Context, shop *models.Shop) error {
	_, err := s.db.ModelContext(ctx, shop).WherePK().Update()
	return err
}

// Product operations
func (s *Store) CreateProduct(ctx context.Context, product *models.Product) error {
	_, err := s.db.ModelContext(ctx, product).Insert()
	return err
}

func (s *Store) GetProductByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	product := &models.Product{ID: id}
	err := s.db.ModelContext(ctx, product).WherePK().Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return product, nil
}

func (s *Store) GetProductsByShopID(ctx context.Context, shopID uuid.UUID, limit, offset int) ([]*models.Product, error) {
	var products []*models.Product
	err := s.db.ModelContext(ctx, &products).
		Where("shop_id = ?", shopID).
		Order("created_at ASC").
		Limit(limit).
		Offset(offset).
		Select()
	return products, err
}

func (s *Store) ListProducts(ctx context.Context, limit, offset int) ([]*models.Product, error) {
	var products []*models.Product
	err := s.db.ModelContext(ctx, &products).
		Order("created_at ASC").
		Limit(limit).
		Offset(offset).
		Select()
	return products, err
}

func (s *Store) SearchProducts(ctx context.Context, query string, limit, offset int) ([]*models.Product, error) {
	var products []*models.Product
	err := s.db.ModelContext(ctx, &products).
		Where("name ILIKE ? OR description ILIKE ?", "%"+query+"%", "%"+query+"%").
		Order("created_at ASC").
		Limit(limit).
		Offset(offset).
		Select()
	return products, err
}

func (s *Store) FilterProductsByCategory(ctx context.Context, category string, limit, offset int) ([]*models.Product, error) {
	var products []*models.Product
	err := s.db.ModelContext(ctx, &products).
		Where("category = ?", category).
		Order("created_at ASC").
		Limit(limit).
		Offset(offset).
		Select()
	return products, err
}

func (s *Store) FilterProductsByPrice(ctx context.Context, minPrice, maxPrice float64, limit, offset int) ([]*models.Product, error) {
	var products []*models.Product
	err := s.db.ModelContext(ctx, &products).
		Where("price BETWEEN ? AND ?", minPrice, maxPrice).
		Order("created_at ASC").
		Limit(limit).
		Offset(offset).
		Select()
	return products, err
}

func (s *Store) UpdateProduct(ctx context.Context, product *models.Product) error {
	_, err := s.db.ModelContext(ctx, product).WherePK().Update()
	return err
}

func (s *Store) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	product := &models.Product{ID: id}
	_, err := s.db.ModelContext(ctx, product).WherePK().Delete()
	return err
}

// Order operations
func (s *Store) CreateOrder(ctx context.Context, order *models.Order) error {
	_, err := s.db.ModelContext(ctx, order).Insert()
	return err
}

func (s *Store) AddOrderItem(ctx context.Context, item *models.OrderItem) error {
	_, err := s.db.ModelContext(ctx, item).Insert()
	return err
}

func (s *Store) GetOrderByID(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	order := &models.Order{ID: id}
	err := s.db.ModelContext(ctx, order).WherePK().Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, errors.New("order not found")
		}
		return nil, err
	}
	return order, nil
}

func (s *Store) GetOrdersByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Order, error) {
	var orders []*models.Order
	err := s.db.ModelContext(ctx, &orders).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Select()
	return orders, err
}

func (s *Store) GetOrderItems(ctx context.Context, orderID uuid.UUID) ([]*models.OrderItem, error) {
	var items []*models.OrderItem
	err := s.db.ModelContext(ctx, &items).
		Where("order_id = ?", orderID).
		Select()
	return items, err
}

func (s *Store) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status models.OrderStatus) (*models.Order, error) {
	order := &models.Order{ID: orderID}
	err := s.db.ModelContext(ctx, order).WherePK().Select()
	if err != nil {
		return nil, err
	}

	order.Status = status
	_, err = s.db.ModelContext(ctx, order).WherePK().Update()
	return order, err
}

// Thêm method để lấy đơn hàng theo shop
func (s *Store) GetOrdersByShopID(ctx context.Context, shopID uuid.UUID, limit, offset int) ([]*models.Order, error) {
	var orders []*models.Order
	err := s.db.ModelContext(ctx, &orders).
		Where("shop_id = ?", shopID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Select()
	return orders, err
}

// Thêm thống kê đơn hàng cho shop
func (s *Store) GetShopOrderStatistics(ctx context.Context, shopID uuid.UUID, startDate, endDate time.Time) (map[string]interface{}, error) {
	var result struct {
		TotalOrders  int     `pg:"total_orders"`
		TotalRevenue float64 `pg:"total_revenue"`
	}

	_, err := s.db.QueryOneContext(ctx, &result, `
		SELECT COUNT(*) as total_orders, SUM(total_amount) as total_revenue
		FROM orders
		WHERE shop_id = ? AND created_at BETWEEN ? AND ?
	`, shopID, startDate, endDate)

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_orders":  result.TotalOrders,
		"total_revenue": result.TotalRevenue,
	}, nil
}

// Transaction support
func (s *Store) RunInTransaction(ctx context.Context, fn func(*pg.Tx) error) error {
	return s.db.RunInTransaction(ctx, fn)
}
