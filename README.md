# Hexa-Go: Implementasi Hexagonal Architecture dengan Go

Proyek ini adalah contoh implementasi **Hexagonal Architecture** (Ports and Adapters) menggunakan bahasa pemrograman Go. Aplikasi ini menyediakan API REST untuk manajemen artikel dan pengguna dengan struktur yang terorganisir, mudah diuji, dan dapat dirawat.

## 📋 Daftar Isi

- [Apa itu Hexagonal Architecture?](#apa-itu-hexagonal-architecture)
- [Struktur Proyek](#struktur-proyek)
- [Komponen Arsitektur](#komponen-arsitektur)
- [Alur Data](#alur-data)
- [Teknologi yang Digunakan](#teknologi-yang-digunakan)
- [Persyaratan](#persyaratan)
- [Instalasi dan Konfigurasi](#instalasi-dan-konfigurasi)
- [Menjalankan Aplikasi](#menjalankan-aplikasi)
- [Struktur API](#struktur-api)

## 🏗️ Apa itu Hexagonal Architecture?

**Hexagonal Architecture** (juga dikenal sebagai **Ports and Adapters**) adalah pola arsitektur yang memisahkan logika bisnis aplikasi dari infrastruktur eksternal. Arsitektur ini disebut "hexagonal" karena dapat digambarkan sebagai hexagon, di mana:

- **Pusat (Core)**: Berisi logika bisnis murni yang tidak bergantung pada teknologi eksternal
- **Ports**: Interface yang mendefinisikan kontrak komunikasi
- **Adapters**: Implementasi konkret yang menghubungkan aplikasi dengan dunia luar

### Keuntungan Hexagonal Architecture

1. **Independensi dari Framework**: Logika bisnis tidak terikat pada framework tertentu
2. **Testabilitas**: Mudah diuji karena dependencies dapat di-mock
3. **Fleksibilitas**: Mudah mengganti teknologi eksternal (database, cache, dll) tanpa mengubah logika bisnis
4. **Maintainability**: Kode lebih terorganisir dan mudah dirawat
5. **Separation of Concerns**: Setiap layer memiliki tanggung jawab yang jelas

### Konsep Dasar

```
┌─────────────────────────────────────────┐
│         Driving Adapters                │
│  (HTTP, CLI, gRPC, Message Queue)       │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│         Application Layer               │
│      (Use Cases / Services)             │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│          Domain Layer                   │
│  (Entities, Business Logic, Ports)      │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│        Driven Adapters                  │
│  (Database, Cache, External APIs)       │
└─────────────────────────────────────────┘
```

## 📁 Struktur Proyek

```
hexa-go/
├── cmd/
│   └── api/
│       └── main.go                    # Entry point aplikasi
├── internal/
│   ├── domain/                        # Domain Layer (Core Business Logic)
│   │   ├── article/
│   │   │   ├── entity.go              # Entity Article
│   │   │   ├── repository.go          # Port (Interface) untuk repository
│   │   │   ├── service.go             # Domain service
│   │   │   └── errors.go              # Domain errors
│   │   └── user/
│   │       ├── entity.go
│   │       ├── repository.go
│   │       ├── service.go
│   │       └── errors.go
│   │   └── media/
│   │       ├── entity.go
│   │       ├── repository.go
│   │       ├── service.go
│   │       ├── storage.go
│   │       └── errors.go
│   │
│   ├── application/                   # Application Layer (Use Cases)
│   │   ├── article/
│   │   │   ├── dto/                   # Data Transfer Objects
│   │   │   │   ├── request.go
│   │   │   │   └── response.go
│   │   │   └── usecase/
│   │   │       ├── create.go
│   │   │       ├── get.go
│   │   │       ├── list.go
│   │   │       ├── update.go
│   │   │       └── delete.go
│   │   └── user/
│   │       ├── dto/                   # Data Transfer Objects
│   │       │   ├── request.go
│   │       │   └── response.go
│   │       └── usecase/
│   │           ├── create.go
│   │           ├── get.go
│   │           ├── list.go
│   │           ├── update.go
│   │           ├── delete.go
│   │           └── login.go
│   │   └── media/
│   │       ├── dto/
│   │       │   ├── request.go
│   │       │   └── response.go
│   │       └── usecase/
│   │           ├── create.go
│   │           ├── get.go
│   │           ├── list.go
│   │           ├── update.go
│   │           └── delete.go
│   │
│   ├── adapters/                      # Adapters Layer
│   │   ├── http/                      # Driving Adapter (HTTP)
│   │   │   ├── article/
│   │   │   │   └── handler.go
│   │   │   ├── user/
│   │   │   │   └── handler.go
│   │   │   ├── media/
│   │   │   │   └── handler.go
│   │   │   ├── response/              # Standard response format
│   │   │   │   └── response.go
│   │   │   ├── router.go
│   │   │   └── middleware.go
│   │   ├── db/                        # Driven Adapter (Database)
│   │   │   ├── article/
│   │   │   │   └── mysql_repo.go
│   │   │   ├── user/
│   │   │   │   └── mysql_repo.go
│   │   │   └── media/
│   │   │       └── mysql_repo.go
│   │   ├── cache/                     # Driven Adapter (Cache)
│   │   │   └── article/
│   │   │       └── redis_cache.go
│   │   ├── storage/                   # Driven Adapter (File Storage)
│   │   │   └── media/
│   │   │       └── local_storage.go
│   │   └── external/                  # Driven Adapter (External Services)
│   │       └── user/
│   │           └── email_sender.go
│   │
│   └── infrastructure/                # Infrastructure Layer
│       ├── config/
│       │   └── config.go              # Konfigurasi aplikasi
│       ├── database/
│       │   ├── mysql.go               # Koneksi MySQL
│       │   └── redis.go               # Koneksi Redis
│       ├── di/                        # Dependency Injection
│       │   ├── container.go
│       │   ├── article/
│       │   │   └── container.go
│       │   ├── user/
│       │   │   └── container.go
│       │   └── media/
│       │       └── container.go
│       └── logger/
│           └── logger.go
├── migration/                         # Database migrations
│   ├── article.sql
│   ├── user.sql
│   └── media.sql
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## 🔧 Komponen Arsitektur

### 1. Domain Layer (`internal/domain/`)

**Domain Layer** adalah inti dari aplikasi yang berisi:

- **Entities**: Objek bisnis dengan logika domain
- **Ports (Interfaces)**: Kontrak yang didefinisikan oleh domain
- **Domain Services**: Logika bisnis yang tidak cocok di entity
- **Domain Errors**: Error spesifik domain

**Contoh: `internal/domain/article/repository.go`**
```go
// Repository adalah port (interface) yang didefinisikan oleh domain
// Domain tidak peduli bagaimana implementasinya
type Repository interface {
    Create(ctx context.Context, article *Article) (*Article, error)
    GetByID(ctx context.Context, id int64) (*Article, error)
    // ... methods lainnya
}
```

**Karakteristik:**
- ✅ Tidak bergantung pada framework atau teknologi eksternal
- ✅ Hanya berisi logika bisnis murni
- ✅ Mendefinisikan kontrak (ports) yang harus dipenuhi oleh adapters

### 2. Application Layer (`internal/application/`)

**Application Layer** berisi use cases yang mengorkestrasikan domain logic:

- **Use Cases**: Setiap use case mewakili satu operasi bisnis
- **DTOs**: Data Transfer Objects untuk komunikasi antar layer

**Contoh: `internal/application/article/usecase/get.go`**
```go
// GetArticleUseCase mengorkestrasikan logika untuk mendapatkan artikel
type GetArticleUseCase struct {
    articleRepo domainarticle.Repository  // Port dari domain
    cache       ArticleSingleCache        // Port untuk cache
}

func (uc *GetArticleUseCase) Execute(ctx context.Context, id int64) (*dto.ArticleResponse, error) {
    // 1. Cek cache dulu
    // 2. Jika tidak ada, ambil dari repository
    // 3. Simpan ke cache
    // 4. Return response
}
```

**Karakteristik:**
- ✅ Mengorkestrasikan domain logic
- ✅ Menggunakan ports (interfaces) dari domain
- ✅ Tidak tahu implementasi konkret dari adapters

### 3. Adapters Layer (`internal/adapters/`)

**Adapters** adalah implementasi konkret yang menghubungkan aplikasi dengan dunia luar:

#### Driving Adapters (Input)
- **HTTP Handler**: Menerima request HTTP dan memanggil use cases
- **CLI**: Command line interface (jika ada)
- **gRPC**: gRPC handlers (jika ada)

**Contoh: `internal/adapters/http/article/handler.go`**
```go
// Handler adalah driving adapter yang menerima HTTP request
type Handler struct {
    getUseCase *article.GetArticleUseCase  // Menggunakan use case
}

func (h *Handler) Get(c *gin.Context) {
    // 1. Parse request
    // 2. Panggil use case
    // 3. Return HTTP response menggunakan standar response
    response, err := h.getUseCase.Execute(ctx, id)
    if err != nil {
        response.ErrorResponseNotFound(c, err.Error())
        return
    }
    response.SuccessResponseOK(c, "Article retrieved successfully", response)
}
```

**Contoh: `internal/adapters/http/response/response.go`**
```go
// StandardResponse adalah struktur response standar untuk semua endpoint
type StandardResponse struct {
    Status  ResponseStatus `json:"status"`  // "success" atau "error"
    Message string         `json:"message"` // Pesan yang menjelaskan hasil
    Data    interface{}    `json:"data,omitempty"` // Data (optional)
}

// Helper functions untuk response
func SuccessResponseOK(c *gin.Context, message string, data interface{}) {
    c.JSON(200, StandardResponse{
        Status:  StatusSuccess,
        Message: message,
        Data:    data,
    })
}

func ErrorResponseNotFound(c *gin.Context, message string) {
    c.JSON(404, StandardResponse{
        Status:  StatusError,
        Message: message,
    })
}
```

#### Driven Adapters (Output)
- **Database Repository**: Implementasi repository untuk MySQL
- **Cache**: Implementasi cache untuk Redis
- **External Services**: Integrasi dengan API eksternal

**Contoh: `internal/adapters/db/article/mysql_repo.go`**
```go
// Repository adalah driven adapter yang mengimplementasikan
// domain article.Repository interface
type Repository struct {
    db *sql.DB
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*article.Article, error) {
    // Implementasi konkret menggunakan MySQL
}
```

**Karakteristik:**
- ✅ Mengimplementasikan ports yang didefinisikan domain
- ✅ Mengetahui detail teknologi eksternal (MySQL, Redis, dll)
- ✅ Dapat diganti tanpa mengubah domain atau application layer

### 4. Infrastructure Layer (`internal/infrastructure/`)

**Infrastructure Layer** menyediakan:

- **Configuration**: Load konfigurasi dari environment
- **Database Connections**: Setup koneksi database
- **Dependency Injection**: Wiring semua dependencies
- **Logging**: Setup logger

**Contoh: `internal/infrastructure/di/article/container.go`**
```go
// NewContainer melakukan dependency injection
func NewContainer(database *sql.DB, redisClient *redis.Client) *Container {
    // 1. Buat repository (driven adapter)
    articleRepo := dbarticle.NewRepository(database)
    
    // 2. Buat cache (driven adapter)
    articleCache := cachearticle.NewRedisCache(redisClient, 5*time.Minute)
    
    // 3. Buat domain service
    articleService := domainarticle.NewService(articleRepo)
    
    // 4. Buat use cases
    getArticleUseCase := apparticle.NewGetArticleUseCase(articleRepo, articleCache)
    
    // 5. Buat handler (driving adapter)
    articleHandler := httparticle.NewHandler(getArticleUseCase, ...)
    
    return &Container{...}
}
```

## 🔄 Alur Data

Berikut adalah alur data ketika user melakukan request untuk mendapatkan artikel:

```
1. HTTP Request
   ↓
2. HTTP Handler (Driving Adapter)
   - Parse request
   - Validasi input
   ↓
3. Use Case (Application Layer)
   - Cek cache
   - Jika tidak ada, panggil repository
   - Simpan ke cache
   ↓
4. Domain Service (Domain Layer)
   - Validasi business rules
   ↓
5. Repository (Driven Adapter)
   - Query ke MySQL
   ↓
6. Response kembali melalui layer yang sama
```

**Contoh Flow:**

```
GET /articles/1
  ↓
ArticleHandler.Get()
  ↓
GetArticleUseCase.Execute()
  ↓
ArticleCache.GetArticle() [Cek cache]
  ↓ (jika tidak ada)
ArticleRepository.GetByID() [Query MySQL]
  ↓
ArticleCache.SetArticle() [Simpan ke cache]
  ↓
Return ArticleResponse
  ↓
HTTP 200 OK dengan JSON response
```

## 🛠️ Teknologi yang Digunakan

- **Go 1.23+**: Bahasa pemrograman
- **Gin**: Web framework untuk HTTP handlers
- **MySQL**: Database untuk persistence
- **Redis**: Cache layer
- **JWT**: Authentication
- **godotenv**: Environment configuration

## 📋 Persyaratan

- Go 1.23 atau lebih tinggi
- MySQL 5.7+ atau MySQL 8.0+
- Redis (opsional, untuk caching)
- Git

## ⚙️ Instalasi dan Konfigurasi

1. **Clone repository**
```bash
git clone <repository-url>
cd hexa-go
```

2. **Install dependencies**
```bash
go mod download
```

3. **Setup environment variables**

Buat file `.env` di root project:
```env
# Server Configuration
SERVER_HOST=localhost
SERVER_PORT=8080
SERVER_DEBUG=true

# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=hexa_go
DB_CHARSET=utf8mb4

# Redis Configuration (Optional)
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT Configuration
JWT_SECRET=your-secret-key-change-in-production
JWT_EXPIRATION=24

# Storage Configuration
STORAGE_BASE_PATH=./storage
STORAGE_BASE_URL=http://localhost:8080
```

4. **Setup Database**

Buat database MySQL:
```sql
CREATE DATABASE hexa_go CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

Jalankan migration SQL untuk membuat tabel:
```bash
# Masuk ke MySQL
mysql -u root -p hexa_go

# Jalankan migration
source migration/user.sql
source migration/article.sql
source migration/media.sql
```

Atau jalankan file SQL secara langsung:
```bash
mysql -u root -p hexa_go < migration/user.sql
mysql -u root -p hexa_go < migration/article.sql
mysql -u root -p hexa_go < migration/media.sql
```

## 🚀 Menjalankan Aplikasi

1. **Jalankan aplikasi menggunakan Go**
```bash
go run cmd/api/main.go
```

Atau menggunakan Make:
```bash
make run
```

2. **Build aplikasi**
```bash
make build-generate
```

3. **Aplikasi akan berjalan di**
```
http://localhost:8080
```

4. **Test aplikasi**
```bash
make test
```

3. **Test API**

Gunakan curl atau Postman untuk test API:

```bash
# Register user
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123"
  }'

# Response:
# {
#   "status": "success",
#   "message": "User registered successfully",
#   "data": { ... }
# }

# Login
curl -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'

# Response:
# {
#   "status": "success",
#   "message": "Login successful",
#   "data": {
#     "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
#     "user": { ... }
#   }
# }

# Create article (dengan token JWT)
curl -X POST http://localhost:8080/api/v1/articles \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "title": "My First Article",
    "content": "This is the content of my article"
  }'

# Response:
# {
#   "status": "success",
#   "message": "Article created successfully",
#   "data": { ... }
# }

# Get article
curl -X GET http://localhost:8080/api/v1/articles/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Response:
# {
#   "status": "success",
#   "message": "Article retrieved successfully",
#   "data": { ... }
# }

# List articles
curl -X GET "http://localhost:8080/api/v1/articles?limit=10&offset=0" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Response:
# {
#   "status": "success",
#   "message": "Articles retrieved successfully",
#   "data": {
#     "articles": [ ... ],
#     "total": 100,
#     "limit": 10,
#     "offset": 0
#   }
# }

# Upload media file
curl -X POST http://localhost:8080/api/v1/media \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "file=@/path/to/your/file.jpg"

# Response:
# {
#   "status": "success",
#   "message": "Media created successfully",
#   "data": {
#     "id": 1,
#     "name": "file.jpg",
#     "path": "2024/01/15/file_1705320000.jpg",
#     "url": "http://localhost:8080/api/v1/media/files/2024/01/15/file_1705320000.jpg",
#     "created_at": "2024-01-15T10:00:00Z",
#     "updated_at": "2024-01-15T10:00:00Z"
#   }
# }

# Get media by ID
curl -X GET http://localhost:8080/api/v1/media/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# List media
curl -X GET "http://localhost:8080/api/v1/media?limit=10&offset=0" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Update media (upload new file)
curl -X PUT http://localhost:8080/api/v1/media/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "file=@/path/to/new/file.jpg"

# Delete media
curl -X DELETE http://localhost:8080/api/v1/media/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## 📡 Struktur API

### User Endpoints

**Public Routes (No Authentication Required):**
- `POST /api/v1/users/register` - Register user baru
- `POST /api/v1/users/login` - Login user

**Protected Routes (Authentication Required):**
- `POST /api/v1/users` - Create user baru
- `GET /api/v1/users` - List users
- `GET /api/v1/users/:id` - Get user by ID
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user

### Article Endpoints

**Protected Routes (Authentication Required):**
- `POST /api/v1/articles` - Create article
- `GET /api/v1/articles` - List articles (with pagination)
- `GET /api/v1/articles/:id` - Get article by ID
- `PUT /api/v1/articles/:id` - Update article
- `DELETE /api/v1/articles/:id` - Delete article

### Media Endpoints

**Public Routes (No Authentication Required):**
- `GET /api/v1/media/files/*` - Access uploaded media files

**Protected Routes (Authentication Required):**
- `POST /api/v1/media` - Upload media file (multipart/form-data, field: `file`)
- `GET /api/v1/media` - List media (with pagination)
- `GET /api/v1/media/:id` - Get media by ID
- `PUT /api/v1/media/:id` - Update media (upload new file, multipart/form-data, field: `file`)
- `DELETE /api/v1/media/:id` - Delete media

**Media Response Format:**
```json
{
  "id": 1,
  "name": "image.jpg",
  "path": "2024/01/15/image_1705320000.jpg",
  "url": "http://localhost:8080/api/v1/media/files/2024/01/15/image_1705320000.jpg",
  "created_at": "2024-01-15T10:00:00Z",
  "updated_at": "2024-01-15T10:00:00Z"
}
```

**Note:**
- Files are stored in directory structure: `YYYY/MM/DD/filename_timestamp.ext`
- Files can be accessed via the `url` field in the response
- The `path` field contains the relative storage path
- When updating media, the old file is automatically deleted

### Health Check

- `GET /health` - Health check endpoint (no authentication required)

## 📦 Format Response Standar

Semua endpoint API menggunakan format response yang standar untuk memastikan konsistensi. Format response terdiri dari tiga field utama:

### Struktur Response

```json
{
  "status": "success" | "error",
  "message": "Pesan yang menjelaskan hasil operasi",
  "data": {} // Optional, hanya ada pada success response
}
```

### Success Response

Response sukses memiliki `status: "success"` dan field `data` yang berisi data yang diminta:

**Contoh: Get User Success (200 OK)**
```json
{
  "status": "success",
  "message": "User retrieved successfully",
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

**Contoh: Create Article Success (201 Created)**
```json
{
  "status": "success",
  "message": "Article created successfully",
  "data": {
    "id": 1,
    "title": "My First Article",
    "content": "This is the content",
    "author_id": 1,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

**Contoh: Delete Success (200 OK)**
```json
{
  "status": "success",
  "message": "Article deleted successfully",
  "data": null
}
```

### Error Response

Response error memiliki `status: "error"` dan tidak memiliki field `data`:

**Contoh: Not Found (404)**
```json
{
  "status": "error",
  "message": "User not found"
}
```

**Contoh: Bad Request (400)**
```json
{
  "status": "error",
  "message": "invalid user id"
}
```

**Contoh: Unauthorized (401)**
```json
{
  "status": "error",
  "message": "authorization header is required"
}
```

**Contoh: Conflict (409)**
```json
{
  "status": "error",
  "message": "email already exists"
}
```

**Contoh: Internal Server Error (500)**
```json
{
  "status": "error",
  "message": "internal server error"
}
```

### Standarisasi Status Code

Aplikasi menggunakan standar HTTP status codes yang konsisten:

#### Success Status Codes
- `200 OK` - Request berhasil (GET, PUT, DELETE)
- `201 Created` - Resource berhasil dibuat (POST)
- `202 Accepted` - Request diterima untuk diproses
- `204 No Content` - Request berhasil tanpa response body

#### Client Error Status Codes
- `400 Bad Request` - Request tidak valid (validation error, invalid parameter)
- `401 Unauthorized` - Tidak terautentikasi (missing/invalid token)
- `403 Forbidden` - Tidak memiliki akses
- `404 Not Found` - Resource tidak ditemukan
- `405 Method Not Allowed` - HTTP method tidak diizinkan
- `409 Conflict` - Konflik dengan state saat ini (e.g., duplicate email)
- `422 Unprocessable Entity` - Request valid tetapi tidak dapat diproses
- `429 Too Many Requests` - Rate limit terlampaui

#### Server Error Status Codes
- `500 Internal Server Error` - Error server internal
- `502 Bad Gateway` - Error dari upstream server
- `503 Service Unavailable` - Service tidak tersedia

### Contoh Penggunaan

**Success Response dengan Data:**
```bash
curl -X GET http://localhost:8080/api/v1/users/1 \
  -H "Authorization: Bearer YOUR_TOKEN"

# Response:
{
  "status": "success",
  "message": "User retrieved successfully",
  "data": { ... }
}
```

**Error Response:**
```bash
curl -X GET http://localhost:8080/api/v1/users/999 \
  -H "Authorization: Bearer YOUR_TOKEN"

# Response (404):
{
  "status": "error",
  "message": "User not found"
}
```

**Validation Error:**
```bash
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{"name": "John"}'

# Response (400):
{
  "status": "error",
  "message": "Key: 'CreateUserRequest.Email' Error:Field validation for 'Email' failed on the 'required' tag"
}
```

## 🎯 Prinsip-Prinsip Hexagonal Architecture dalam Proyek Ini

1. **Dependency Inversion**: Domain layer tidak bergantung pada adapters, sebaliknya adapters bergantung pada domain
2. **Interface Segregation**: Setiap port (interface) memiliki tanggung jawab yang spesifik
3. **Single Responsibility**: Setiap layer dan komponen memiliki satu tanggung jawab
4. **Open/Closed Principle**: Mudah menambah adapter baru tanpa mengubah domain logic

## 📚 Referensi

- [Hexagonal Architecture by Alistair Cockburn](https://alistair.cockburn.us/hexagonal-architecture/)
- [Clean Architecture by Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

## 📝 License

MIT License

---

**Dibuat dengan ❤️ menggunakan Go dan Hexagonal Architecture**

