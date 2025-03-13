package models

import (
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/google/uuid"
)

type UserRole string

const (
	RoleAdmin  UserRole = "admin"
	RoleSeller UserRole = "seller"
	RoleBuyer  UserRole = "buyer"
)

type OrderStatus string

const (
	StatusPending   OrderStatus = "pending"
	StatusPaid      OrderStatus = "paid"
	StatusShipped   OrderStatus = "shipped"
	StatusDelivered OrderStatus = "delivered"
	StatusCanceled  OrderStatus = "canceled"
)

type User struct {
	ID           uuid.UUID `pg:"id,type:uuid,pk,default:gen_random_uuid()"`
	Username     string    `pg:"username,unique,notnull"`
	Email        string    `pg:"email,unique,notnull"`
	PasswordHash string    `pg:"password_hash,notnull"`
	Role         UserRole  `pg:"role,notnull,type:user_role,default:'buyer'"`
	IsBanned     bool      `pg:"is_banned,notnull,default:false"`
	CreatedAt    time.Time `pg:"created_at,notnull,default:now()"`
	UpdatedAt    time.Time `pg:"updated_at,notnull,default:now()"`
	// Relations
	Shops  []*Shop  `pg:"rel:has-many"`
	Orders []*Order `pg:"rel:has-many"`
}

type Shop struct {
	ID          uuid.UUID `pg:"id,type:uuid,pk,default:gen_random_uuid()"`
	UserID      uuid.UUID `pg:"user_id,type:uuid,notnull"`
	Name        string    `pg:"name,unique,notnull"`
	Description string    `pg:"description"`
	LogoURL     string    `pg:"logo_url"`
	CreatedAt   time.Time `pg:"created_at,notnull,default:now()"`
	UpdatedAt   time.Time `pg:"updated_at,notnull,default:now()"`
	// Relations
	User     *User      `pg:"rel:belongs-to"`
	Products []*Product `pg:"rel:has-many"`
}

type Product struct {
	ID          uuid.UUID `pg:"id,type:uuid,pk,default:gen_random_uuid()"`
	ShopID      uuid.UUID `pg:"shop_id,type:uuid,notnull"`
	Name        string    `pg:"name,notnull"`
	Description string    `pg:"description"`
	Price       float64   `pg:"price,notnull"`
	Stock       int32     `pg:"stock,notnull,default:0"`
	Category    string    `pg:"category,notnull"`
	ImageURLs   []string  `pg:"image_urls,array"`
	CreatedAt   time.Time `pg:"created_at,notnull,default:now()"`
	UpdatedAt   time.Time `pg:"updated_at,notnull,default:now()"`
	// Relations
	Shop       *Shop        `pg:"rel:belongs-to"`
	OrderItems []*OrderItem `pg:"rel:has-many"`
}

type Order struct {
	ID              uuid.UUID   `pg:"id,type:uuid,pk,default:gen_random_uuid()"`
	UserID          uuid.UUID   `pg:"user_id,type:uuid,notnull"`
	ShopID          uuid.UUID   `pg:"shop_id,type:uuid,notnull"` // Thêm ShopID
	TotalAmount     float64     `pg:"total_amount,notnull"`
	Status          OrderStatus `pg:"status,notnull,type:order_status,default:'pending'"`
	ShippingAddress string      `pg:"shipping_address,notnull"`
	CreatedAt       time.Time   `pg:"created_at,notnull,default:now()"`
	UpdatedAt       time.Time   `pg:"updated_at,notnull,default:now()"`
	// Relations
	User       *User        `pg:"rel:belongs-to"`
	Shop       *Shop        `pg:"rel:belongs-to"` // Thêm quan hệ với Shop
	OrderItems []*OrderItem `pg:"rel:has-many"`
}

type OrderItem struct {
	ID              uuid.UUID `pg:"id,type:uuid,pk,default:gen_random_uuid()"`
	OrderID         uuid.UUID `pg:"order_id,type:uuid,notnull"`
	ProductID       uuid.UUID `pg:"product_id,type:uuid,notnull"`
	Quantity        int32     `pg:"quantity,notnull"`
	PriceAtPurchase float64   `pg:"price_at_purchase,notnull"`
	CreatedAt       time.Time `pg:"created_at,notnull,default:now()"`
	// Relations
	Order   *Order   `pg:"rel:belongs-to"`
	Product *Product `pg:"rel:belongs-to"`
}

// CreateSchema creates database schema for all models
func CreateSchema(db *pg.DB) error {
	models := []interface{}{
		(*User)(nil),
		(*Shop)(nil),
		(*Product)(nil),
		(*Order)(nil),
		(*OrderItem)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
