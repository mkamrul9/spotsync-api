# 🚗 SpotSync API — Smart Parking & EV Reservation System

SpotSync is a centralized backend REST API designed for busy airports and malls to manage parking zones, specifically handling the high-demand reservation of limited EV charging spots.

**Live API URL:** [https://spotsync-api.onrender.com](https://spotsync-api.onrender.com)  

---

## 🛠️ Technology Stack

| Technology | Package | Note |
|---|---|---|
| **Go (Golang)** | — | Version 1.22+ |
| **Echo** | `github.com/labstack/echo/v4` | High-performance web framework |
| **GORM** | `gorm.io/gorm` | ORM with PostgreSQL driver |
| **PostgreSQL** | NeonDB (cloud) | Relational database |
| **Validator** | `github.com/go-playground/validator/v10` | Struct validation integrated with Echo |
| **JWT** | `github.com/golang-jwt/jwt/v5` | Token generation & verification |
| **bcrypt** | `golang.org/x/crypto/bcrypt` | Password hashing (cost 12) |
| **godotenv** | `github.com/joho/godotenv` | Environment variable loading |

---

## 🏛️ Architecture — Clean Architecture (Strict)

Handlers must **NOT** talk to the database directly. Each layer has one responsibility:

```
Handler → Service → Repository → Database
   ↑           ↑
  DTO        Models
```

| Layer | Directory | Responsibility |
|---|---|---|
| **DTO** | `dto/` | Request payloads & response structures. GORM models never exposed to API. |
| **Handler** | `handler/` | Binds/validates DTOs, extracts JWT claims, calls Service, returns JSON. |
| **Service** | `service/` | Business logic: password hashing, JWT generation, capacity enforcement. |
| **Repository** | `repository/` | All GORM database operations (CRUD, Transactions, Row Locks). |
| **Models** | `models/` | GORM structs representing database tables. |

### ⚡ Concurrency Solution — EV Spot Bottleneck

To prevent race conditions (two drivers booking the last spot simultaneously), the `CreateReservation` method uses a **GORM Database Transaction** with **Row-Level Locking (`SELECT ... FOR UPDATE`)**:

```go
db.Transaction(func(tx *gorm.DB) error {
    // 1. Lock the zone row
    tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&zone, zoneID)
    // 2. Count active reservations
    // 3. Reject if at capacity (409 Conflict)
    // 4. Otherwise, create reservation atomically
})
```

---

## 👥 User Roles & Permissions

| Role | Permissions |
|---|---|
| **driver** | Register/login · View zones · Reserve a spot · View/cancel own reservations |
| **admin** | All driver permissions · Create parking zones · View all reservations |

---

## 🚀 Local Development Setup

### Prerequisites

- Go 1.22+ installed
- PostgreSQL database (local, NeonDB, or Supabase)

### 1. Clone the repository

```bash
git clone https://github.com/mkamrul9/spotsync-api.git
cd spotsync-api
```

### 2. Create your `.env` file

```env
PORT=8080
DB_URL=postgres://user:password@host:port/dbname?sslmode=require
JWT_SECRET=your_super_secret_jwt_key
```

### 3. Install dependencies & run

```bash
go mod tidy
go run main.go
```

> **Note:** GORM `AutoMigrate` runs automatically on startup — all database tables (`users`, `parking_zones`, `reservations`) are created for you.

---

## 🌐 API Endpoints

### 🔓 Public Routes

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/` | Welcome message |
| `GET` | `/health` | System health check |
| `POST` | `/api/v1/auth/register` | Register a new user |
| `POST` | `/api/v1/auth/login` | Login and receive JWT token |
| `GET` | `/api/v1/zones` | Get all parking zones (with live `available_spots`) |
| `GET` | `/api/v1/zones/:id` | Get a single parking zone by ID |

### 🔐 Authenticated Routes (Requires `Authorization: Bearer <token>`)

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/api/v1/reservations` | Reserve a parking spot |
| `GET` | `/api/v1/reservations/my-reservations` | View your own reservations |
| `DELETE` | `/api/v1/reservations/:id` | Cancel your reservation |

### 🛡️ Admin-Only Routes (Requires JWT + `admin` role)

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/api/v1/zones` | Create a new parking zone |
| `PUT` | `/api/v1/zones/:id` | Update an existing parking zone (partial update) |
| `DELETE` | `/api/v1/zones/:id` | Delete a parking zone |
| `GET` | `/api/v1/reservations` | View all reservations in the system |

---

## 📋 Response Format

All responses follow a consistent structure:

**Success:**

```json
{
  "success": true,
  "message": "Operation description",
  "data": {}
}
```

**Error:**

```json
{
  "success": false,
  "message": "Error description",
  "errors": "Error details"
}
```

### HTTP Status Codes

| Code | Meaning |
|---|---|
| `200` | OK — Successful GET/DELETE |
| `201` | Created — Successful POST |
| `400` | Bad Request — Validation errors |
| `401` | Unauthorized — Missing/invalid JWT |
| `403` | Forbidden — Insufficient role |
| `404` | Not Found — Resource does not exist |
| `409` | Conflict — Zone is full / duplicate resource |
| `500` | Internal Server Error |
