# Savannah Store

Built in Go, with PostgreSQL and Auth0 for authentication
Deployment of the API is done on Contabo VPS

---

## Table of Contents

- [Features](#features)
- [Architecture](#architecture)
- [Getting Started](#getting-started)
- [API Endpoints](#api-endpoints)
  - [Authentication](#authentication)
  - [Products](#products)
  - [Categories](#categories)
  - [Cart](#cart)
  - [Orders](#orders)
- [Data Models](#data-models)
- [Authentication & Security](#authentication--security)

---

## Features

- **User Authentication** (Auth0, JWT, OIDC)
- **Product & Category Management**
- **Shopping Cart** (add, remove, update items)
- **Order Processing** (create, list, update status)
- **Role-based Access Control** (Customer, Admin, Super Admin)
- **RESTful API** with [Swagger UI](http://localhost:8080/swagger/index.html)
- **PostgreSQL** database integration
- **Environment-based configuration**
- **Automated Testing** (integration & e2e)

---

## Architecture

- **cmd/server**: Application entrypoint
- **internal/api**: API route registration
- **internal/handlers**: HTTP request handlers
- **internal/model**: Data models (User, Product, Cart, Order, etc.)
- **internal/repocitory**: Database repositories (CRUD logic)
- **internal/middleware**: Auth & request middleware
- **internal/service**: External integrations (email, SMS)
- **internal/database**: DB connection logic
- **migrations**: SQL migration scripts

---

## Getting Started

### Prerequisites

- Go 1.24+
- Docker

### Setup

1. Clone the repository:
   ```sh
   git clone https://github.com/Oj-washingtone/savannah-store.git
   cd savannah-store
   ```
2. Configure environment variables in `.env` (see sample in docs).
3. Run
   ```sh
    docker compose up
   ```
4. Access Swagger UI at [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

---

## API Endpoints

### Authentication

- `POST /api/auth/login` — Login via Auth0
- `GET /api/auth/auth0/callback` — Auth0 callback

#### Auth Step 1:

step 1 Open the Login URL

Open the following URL in your browser:

http://savanna.apis.linxs.co.ke/api/auth/login or http://localhost:8080/api/auth/login

This will redirect you to the Auth0 login page
Enter credentials and login

After successful login, Auth0 will redirect you back to the callback URL configured in the application (`/auth/auth0/callback`)
and be redirected to the browser page having your token

```json
{
  "user": {
    "id": "user-uuid",
    "name": "John Doe",
    "email": "johndoe@example.com",
    "role": "customer",
    "auth0Id": "auth0|123456789"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Products

- `POST /api/products/create` — Add new product
- `GET /api/products/:id` — Get product by ID
- `GET /api/products` — List products (pagination)
- `PATCH /api/products/:id` — Update product
- `DELETE /api/products/:id` — Delete product

### Categories

- `POST /api/products/categories/create` — Add category
- `GET /api/products/categories` — List categories
- `PATCH /api/products/categories/:id` — Update category

### Cart

- `POST /api/cart/create` — Add item to cart
- `DELETE /api/cart/remove/:id` — Remove item from cart
- `GET /api/cart` — List cart items
- `PATCH /api/cart/update/quantity/:id` — Update item quantity

### Orders

- `POST /api/orders/create` — Create order (requires authentication)
- `GET /api/orders` — List all orders

---

## Data Models

### User

```go
Auth0Id   string
Name      string
Email     string
Phone     string
Role      string // customer, admin, super_admin
```

### Product

```go
CategoryID  uuid.UUID
Name        string
Description string
Price       int64
Stock       int
```

### ProductCategory

```go
Name     string
ParentId *uuid.UUID
```

### Cart & CartItem

```go
Cart:  UserId uuid.UUID
CartItem: CartId, ProductId, Quantity, Price
```

### Orders & OrderItems

```go
Orders:  UserID, Status, Total, Paid
OrderItems: OrderID, ProductID, Quantity, Price
```

---

### .env file

THe following are the fields you need in .env while running the application

```
PORT=8080


DB_USER=
DB_PASSWORD=
DB_NAME=
DB_HOST=
DB_PORT=

# Auth 0

AUTH0_DOMAIN=
AUTH0_CLIENT_ID=
AUTH0_CLIENT_SECRET=
AUTH0_CALLBACK_URL=


# Email
RESEND_KEY =
ADMIN_EMAIL =

# SMS
AFRICASTALKING_API_KEY =
AFRICASTALKING_URL =
AFRICASTALKING_USERNAME =

```
