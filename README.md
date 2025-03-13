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

- `POST /api/auth/register`: Register a new user
- `POST /api/auth/login`: Login and get JWT token

### User Endpoints

- `GET /api/profile`: Get current user profile

### Shop Endpoints

- `POST /api/shops`: Create a new shop
- `GET /api/shops`: List all shops
- `GET /api/shops/user`: List current user's shops
- `GET /api/shops/:id`: Get shop details
- `PUT /api/shops/:id`: Update shop details
- `GET /api/shops/search`: Search shops

### Product Endpoints

- `GET /api/products`: List all products
- `GET /api/products/:id`: Get product details
- `GET /api/products/search`: Search products
- `GET /api/products/category`: Filter products by category
- `GET /api/products/price`: Filter products by price range
- `GET /api/shops/:id/products`: List products in a shop

### Order Endpoints

- `POST /api/orders`: Create a new order
- `GET /api/orders`: List current user's orders
- `GET /api/orders/:id`: Get order details

### Seller Endpoints

- `POST /api/seller/products`: Create a new product
- `PUT /api/seller/products/:id`: Update product details
- `DELETE /api/seller/products/:id`: Delete a product

### Admin Endpoints

- `GET /api/admin/users`: List all users
- `POST /api/admin/users/:id/ban`: Ban a user
- `POST /api/admin/users/:id/unban`: Unban a user
- `PUT /api/admin/orders/:id/status`: Update order status

## Default Accounts

After running `make seed-go-pg`, you can use these accounts:

- **Admin**: admin@example.com / admin123
- **Seller**: seller1@example.com / seller123
- **Buyer**: buyer1@example.com / buyer1

## License

This project is licensed under the MIT License.
