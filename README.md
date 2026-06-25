# 🚗 SpotSync API - Smart Parking & EV Reservation System

SpotSync is a centralized backend platform designed for busy airports and malls to manage parking zones, specifically handling the high-demand reservation of limited EV charging spots. 

**Live API URL:** `https://spotsync-api.onrender.com` *(Replace with your actual URL)*  
**Interview Video:** `https://youtu.be/your-video-link` *(Replace with your actual URL)*

## 🛠️ Technology Stack
* **Language:** Go (Golang) v1.22+
* **Framework:** Echo (v4)
* **ORM:** GORM
* **Database:** PostgreSQL (NeonDB / Supabase)
* **Security:** JWT (golang-jwt/v5) & bcrypt

## 🏛️ Architecture (Clean Architecture)
This project strictly follows Clean Architecture principles to separate concerns and ensure maintainability:
1. **Handler Layer (`/handler`):** Parses HTTP requests, validates DTOs, and formats JSON responses. Has no knowledge of the database.
2. **Service Layer (`/service`):** Contains business logic (password hashing, capacity calculations, JWT generation).
3. **Repository Layer (`/repository`):** Handles all GORM database operations.
4. **Data Transfer Objects (`/dto`):** Defines strict request/response structures to prevent mass-assignment vulnerabilities.

**Concurrency Solution:** To prevent the "EV Spot Bottleneck" (overbooking a full zone), the `CreateReservation` repository method utilizes a **GORM Database Transaction** combined with **Row-Level Locking (`FOR UPDATE`)**. This ensures atomicity when checking capacity and inserting the reservation.

## 🚀 Setup & Local Development

### Prerequisites
* Go 1.22+ installed
* PostgreSQL running locally or via cloud (NeonDB/Supabase)

### 1. Clone the repository
```bash
git clone https://github.com/mkamrul9/spotsync-api.git
cd spotsync-api
```

### 2. Environment Variables
Create a `.env` file in the root directory:

```env
PORT=8080
DB_URL=postgres://user:password@host:port/dbname?sslmode=require
JWT_SECRET=your_super_secret_jwt_key
```

### 3. Run the application
```bash
go mod tidy
go run main.go
```
*Note: The application will automatically run GORM AutoMigrate to set up your database tables.*

## 🌐 API Endpoints

### Public Routes
* `GET /health` - System health check
* `POST /api/v1/auth/register` - Register a new user (driver/admin)
* `POST /api/v1/auth/login` - Authenticate and receive JWT
* `GET /api/v1/zones` - View all parking zones and available spots

### Authenticated Routes (Requires JWT Bearer Token)
* `POST /api/v1/reservations` - Reserve a parking spot
* `GET /api/v1/reservations/my-reservations` - View your reservations
* `DELETE /api/v1/reservations/:id` - Cancel your reservation

### Admin-Only Routes (Requires JWT + Admin Role)
* `POST /api/v1/zones` - Create a new parking zone
* `GET /api/v1/reservations` - View all system reservations
