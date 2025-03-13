package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"github.com/qhh/prjEcom/pkg/models"
	"github.com/qhh/prjEcom/pkg/utils"
)

// SeedDatabase seeds the database with test data
func SeedDatabase(db *pg.DB) error {
	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	ctx := context.Background()

	// Seed admin user
	adminPassword, _ := utils.HashPassword("admin123")
	adminUser := &models.User{
		ID:           uuid.New(),
		Username:     "admin",
		Email:        "admin@example.com",
		PasswordHash: adminPassword,
		Role:         models.RoleAdmin,
	}

	_, err := db.ModelContext(ctx, adminUser).Insert()
	if err != nil {
		log.Printf("Error creating admin: %v", err)
	} else {
		log.Printf("Admin created: %s", adminUser.Username)
	}

	// Seed sellers
	sellers := []struct {
		username string
		email    string
		password string
	}{
		{"seller1", "seller1@example.com", "seller123"},
		{"seller2", "seller2@example.com", "seller123"},
	}

	for _, s := range sellers {
		sellerPassword, _ := utils.HashPassword(s.password)
		sellerUser := &models.User{
			ID:           uuid.New(),
			Username:     s.username,
			Email:        s.email,
			PasswordHash: sellerPassword,
			Role:         models.RoleSeller,
		}

		_, err := db.ModelContext(ctx, sellerUser).Insert()
		if err != nil {
			log.Printf("Error creating seller %s: %v", s.username, err)
			continue
		}
		log.Printf("Seller created: %s", sellerUser.Username)

		// Create shops for sellers
		shop := &models.Shop{
			ID:          uuid.New(),
			UserID:      sellerUser.ID,
			Name:        fmt.Sprintf("%s's Shop", sellerUser.Username),
			Description: fmt.Sprintf("Welcome to %s's Shop!", sellerUser.Username),
			// LogoUrl:     "https://via.placeholder.com/150",
		}

		_, err = db.ModelContext(ctx, shop).Insert()
		if err != nil {
			log.Printf("Error creating shop for %s: %v", sellerUser.Username, err)
			continue
		}
		log.Printf("Shop created: %s", shop.Name)

		// Add products to shops
		categories := []string{"Electronics", "Clothing", "Books", "Home", "Sports"}
		for i := 0; i < 5; i++ {
			price := float64(10 + rand.Intn(990))
			stock := int32(5 + rand.Intn(100))
			category := categories[rand.Intn(len(categories))]

			product := &models.Product{
				ID:          uuid.New(),
				ShopID:      shop.ID,
				Name:        fmt.Sprintf("Product %d from %s", i+1, sellerUser.Username),
				Description: fmt.Sprintf("This is product %d from %s's shop", i+1, sellerUser.Username),
				Price:       price,
				Stock:       stock,
				Category:    category,
				// ImageUrls:   []string{fmt.Sprintf("https://via.placeholder.com/300?text=Product%d", i+1)},
			}

			_, err := db.ModelContext(ctx, product).Insert()
			if err != nil {
				log.Printf("Error creating product: %v", err)
				continue
			}
			log.Printf("Product created: %s", product.Name)
		}
	}

	// Seed buyers
	for i := 1; i <= 5; i++ {
		buyerPassword, _ := utils.HashPassword(fmt.Sprintf("buyer%d", i))
		buyerUser := &models.User{
			ID:           uuid.New(),
			Username:     fmt.Sprintf("buyer%d", i),
			Email:        fmt.Sprintf("buyer%d@example.com", i),
			PasswordHash: buyerPassword,
			Role:         models.RoleBuyer,
		}

		_, err := db.ModelContext(ctx, buyerUser).Insert()
		if err != nil {
			log.Printf("Error creating buyer %d: %v", i, err)
			continue
		}
		log.Printf("Buyer created: %s", buyerUser.Username)
	}

	log.Println("Seeding completed.")
	return nil
}
