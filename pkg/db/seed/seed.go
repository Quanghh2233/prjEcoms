package main

// import (
// 	"context"
// 	"database/sql"
// 	"fmt"
// 	"log"
// 	"math/rand"

// 	_ "github.com/lib/pq"
// 	"github.com/qhh/prjEcom/pkg/db/sqlc"
// 	"github.com/qhh/prjEcom/pkg/utils"
// )

// const (
// 	dbDriver = "postgres"
// 	dbSource = "postgresql://qhh:2203@localhost:5432/ecommerce?sslmode=disable"
// )

// func main() {
// 	db, err := sql.Open(dbDriver, dbSource)
// 	if err != nil {
// 		log.Fatalf("Failed to connect to database: %v", err)
// 	}
// 	defer db.Close()

// 	err = db.Ping()
// 	if err != nil {
// 		log.Fatalf("Failed to ping database: %v", err)
// 	}

// 	queries := sqlc.New(db)
// 	ctx := context.Background()

// 	// Seed admin user
// 	adminPassword, _ := utils.HashPassword("admin123")
// 	admin, err := queries.CreateUser(ctx, sqlc.CreateUserParams{
// 		Username:     "admin",
// 		Email:        "admin@example.com",
// 		PasswordHash: adminPassword,
// 		Role:         "admin",
// 	})
// 	if err != nil {
// 		log.Printf("Error creating admin: %v", err)
// 	} else {
// 		log.Printf("Admin created: %s", admin.Username)
// 	}

// 	// Seed sellers
// 	sellers := []struct {
// 		username string
// 		email    string
// 		password string
// 	}{
// 		{"seller1", "seller1@example.com", "seller123"},
// 		{"seller2", "seller2@example.com", "seller123"},
// 	}

// 	for _, s := range sellers {
// 		sellerPassword, _ := utils.HashPassword(s.password)
// 		seller, err := queries.CreateUser(ctx, sqlc.CreateUserParams{
// 			Username:     s.username,
// 			Email:        s.email,
// 			PasswordHash: sellerPassword,
// 			Role:         "seller",
// 		})
// 		if err != nil {
// 			log.Printf("Error creating seller %s: %v", s.username, err)
// 			continue
// 		}
// 		log.Printf("Seller created: %s", seller.Username)

// 		// Create shops for sellers
// 		shop, err := queries.CreateShop(ctx, sqlc.CreateShopParams{
// 			UserID:      seller.ID,
// 			Name:        fmt.Sprintf("%s's Shop", seller.Username),
// 			Description: fmt.Sprintf("Welcome to %s's Shop!", seller.Username),
// 			LogoUrl:     sql.NullString{String: "https://via.placeholder.com/150", Valid: true},
// 		})
// 		if err != nil {
// 			log.Printf("Error creating shop for %s: %v", seller.Username, err)
// 			continue
// 		}
// 		log.Printf("Shop created: %s", shop.Name)

// 		// Add products to shops
// 		categories := []string{"Electronics", "Clothing", "Books", "Home", "Sports"}
// 		for i := 0; i < 5; i++ {
// 			price := float64(10 + rand.Intn(990))
// 			stock := int32(5 + rand.Intn(100))
// 			category := categories[rand.Intn(len(categories))]

// 			product, err := queries.CreateProduct(ctx, sqlc.CreateProductParams{
// 				ShopID:      shop.ID,
// 				Name:        fmt.Sprintf("Product %d from %s", i+1, seller.Username),
// 				Description: fmt.Sprintf("This is product %d from %s's shop", i+1, seller.Username),
// 				Price:       float64(price),
// 				Stock:       stock,
// 				Category:    category,
// 				ImageUrls:   []string{fmt.Sprintf("https://via.placeholder.com/300?text=Product%d", i+1)},
// 			})
// 			if err != nil {
// 				log.Printf("Error creating product: %v", err)
// 				continue
// 			}
// 			log.Printf("Product created: %s", product.Name)
// 		}
// 	}

// 	// Seed buyers
// 	for i := 1; i <= 5; i++ {
// 		buyerPassword, _ := utils.HashPassword(fmt.Sprintf("buyer%d", i))
// 		buyer, err := queries.CreateUser(ctx, sqlc.CreateUserParams{
// 			Username:     fmt.Sprintf("buyer%d", i),
// 			Email:        fmt.Sprintf("buyer%d@example.com", i),
// 			PasswordHash: buyerPassword,
// 			Role:         "buyer",
// 		})
// 		if err != nil {
// 			log.Printf("Error creating buyer %d: %v", i, err)
// 			continue
// 		}
// 		log.Printf("Buyer created: %s", buyer.Username)
// 	}

// 	log.Println("Seeding completed.")
// }
