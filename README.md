# E-commerce API

This is a complete e-commerce API built with Go, Gin, PostgreSQL, and go-pg.

## Features

- **User Authentication**: Register, login, and profile management
- **Shop Management**: Create and manage shops
- **Product Management**: CRUD operations for products
- **Order Processing**: Create orders, view order history
- **Admin Dashboard**: User management, shop oversight
- **Seller Dashboard**: Product management, order fulfillment
- **Search and Filtering**: Find products by name, category, price

## Role System

- **Buyer**: Default role for new users
- **Seller**: Users who own at least one shop
- **Admin**: System administrators with full access

## Technologies Used

- Go 1.20
- Gin Web Framework
- go-pg ORM for PostgreSQL
- PostgreSQL 14+
- Docker and Docker Compose
- JWT for Authentication

## Prerequisites

- Go 1.20+
- PostgreSQL 14+
- Docker and Docker Compose

## Getting Started

1. Clone the repository:

```sh
git clone https://github.com/yourusername/ecommerce-api.git
cd ecommerce-api
```

2. Copy the example environment file:

```sh
cp .env.example .env
```

3. Edit `.env` to set your own database credentials and JWT secret.

4. Start the services using Docker Compose:

```sh
make up
```

5. Initialize the database schema:

```sh
make schema
```

6. Seed the database with initial data:

```sh
make seed-go-pg
```

## Development

### Available Commands

- `make up`: Start all services
- `make down`: Stop all services
- `make build`: Rebuild the services
- `make schema`: Initialize database schema
- `make seed-go-pg`: Seed the database with initial data using go-pg
- `make start`: Start the API locally (without Docker)
- `make test`: Run tests

## API Documentation

### Authentication Endpoints

#### Register a new user

- **URL**: `POST /api/auth/register`
- **Request Body**:
```json
{
  "username": "newuser",
  "email": "newuser@example.com",
  "password": "securepassword123"
}
```
- **Response**: User object with ID, username, email, role

#### Login

- **URL**: `POST /api/auth/login`
- **Request Body**:
```json
{
  "username": "existinguser",
  "password": "mypassword"
}
```
- **Response**: User info and JWT token

### User Endpoints

#### Get current user profile

- **URL**: `GET /api/profile`
- **Headers**: Authorization: Bearer {token}
- **Response**: User profile data

### Shop Endpoints

#### Create a new shop

- **URL**: `POST /api/shops`
- **Headers**: Authorization: Bearer {token}
- **Request Body**:
```json
{
  "name": "My Awesome Shop",
  "description": "Selling the best products",
  "logo_url": "https://example.com/logo.png"
}
```
- **Response**: Created shop object

#### List all shops

- **URL**: `GET /api/shops?limit=10&offset=0`
- **Headers**: Authorization: Bearer {token}
- **Response**: Array of shop objects

#### List current user's shops

- **URL**: `GET /api/shops/user`
- **Headers**: Authorization: Bearer {token}
- **Response**: Array of shop objects owned by current user

#### Get shop details

- **URL**: `GET /api/shops/{shop_id}`
- **Headers**: Authorization: Bearer {token}
- **Response**: Shop object with details

#### Update shop details

- **URL**: `PUT /api/shops/{shop_id}`
- **Headers**: Authorization: Bearer {token}
- **Request Body**:
```json
{
  "name": "Updated Shop Name",
  "description": "Updated shop description",
  "logo_url": "https://example.com/updated-logo.png"
}
```
- **Response**: Updated shop object

#### Search shops

- **URL**: `GET /api/shops/search?q=keyword&limit=10&offset=0`
- **Headers**: Authorization: Bearer {token}
- **Response**: Array of matching shop objects

### Product Endpoints

#### List all products

- **URL**: `GET /api/products?limit=10&offset=0`
- **Headers**: Authorization: Bearer {token}
- **Response**: Array of product objects

#### Get product details

- **URL**: `GET /api/products/{product_id}`
- **Headers**: Authorization: Bearer {token}
- **Response**: Product object with details

#### Search products

- **URL**: `GET /api/products/search?q=keyword&limit=10&offset=0`
- **Headers**: Authorization: Bearer {token}
- **Response**: Array of matching product objects

#### Filter products by category

- **URL**: `GET /api/products/category?category=Electronics&limit=10&offset=0`
- **Headers**: Authorization: Bearer {token}
- **Response**: Array of filtered product objects

#### Filter products by price range

- **URL**: `GET /api/products/price?min_price=10&max_price=100&limit=10&offset=0`
- **Headers**: Authorization: Bearer {token}
- **Response**: Array of filtered product objects

#### List products in a shop

- **URL**: `GET /api/shops/{shop_id}/products?limit=10&offset=0`
- **Headers**: Authorization: Bearer {token}
- **Response**: Array of products for a specific shop

### Order Endpoints

#### Create a new order

- **URL**: `POST /api/orders`
- **Headers**: Authorization: Bearer {token}
- **Request Body**:
```json
{
  "shop_id": "shop-uuid-here",
  "shipping_address": "123 Main St, City, Country",
  "items": [
    {
      "product_id": "product-uuid-here",
      "quantity": 2,
      "price": 19.99
    },
    {
      "product_id": "another-product-uuid",
      "quantity": 1,
      "price": 29.99
    }
  ]
}
```
- **Response**: Order creation confirmation

#### List current user's orders

- **URL**: `GET /api/orders?limit=10&offset=0`
- **Headers**: Authorization: Bearer {token}
- **Response**: Array of order objects

#### Get order details

- **URL**: `GET /api/orders/{order_id}`
- **Headers**: Authorization: Bearer {token}
- **Response**: Order object with items

### Seller Endpoints (require seller role)

#### Create a new product

- **URL**: `POST /api/seller/products`
- **Headers**: Authorization: Bearer {token}
- **Request Body**:
```json
{
  "shop_id": "shop-uuid-here",
  "name": "New Product",
  "description": "This is a great product",
  "price": 49.99,
  "stock": 100,
  "category": "Electronics",
  "image_urls": [
    "https://example.com/image1.jpg",
    "https://example.com/image2.jpg"
  ]
}
```
- **Response**: Created product object

#### Update product details

- **URL**: `PUT /api/seller/products/{product_id}`
- **Headers**: Authorization: Bearer {token}
- **Request Body**:
```json
{
  "shop_id": "shop-uuid-here",
  "name": "Updated Product Name",
  "description": "Updated description",
  "price": 59.99,
  "stock": 80,
  "category": "Electronics",
  "image_urls": [
    "https://example.com/updated-image1.jpg"
  ]
}
```
- **Response**: Updated product object

#### Delete a product

- **URL**: `DELETE /api/seller/products/{product_id}`
- **Headers**: Authorization: Bearer {token}
- **Response**: Success message

### Admin Endpoints (require admin role)

#### List all users

- **URL**: `GET /api/admin/users?limit=10&offset=0`
- **Headers**: Authorization: Bearer {token}
- **Response**: Array of user objects

#### Ban a user

- **URL**: `POST /api/admin/users/{user_id}/ban`
- **Headers**: Authorization: Bearer {token}
- **Response**: Updated user object

#### Unban a user

- **URL**: `POST /api/admin/users/{user_id}/unban`
- **Headers**: Authorization: Bearer {token}
- **Response**: Updated user object

#### Update order status

- **URL**: `PUT /api/admin/orders/{order_id}/status`
- **Headers**: Authorization: Bearer {token}
- **Request Body**:
```json
{
  "status": "shipped"
}
```
- **Response**: Updated order object

## Default Accounts

After running `make seed-go-pg`, you can use these accounts:

- **Admin**: admin@example.com / admin123
- **Seller**: seller1@example.com / seller123
- **Buyer**: buyer1@example.com / buyer1

## License

This project is licensed under the MIT License.
