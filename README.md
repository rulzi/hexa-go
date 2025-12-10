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
│   │
│   ├── application/                   # Application Layer (Use Cases)
│   │   ├── article/
│   │   │   ├── dto.go                 # Data Transfer Objects
│   │   │   ├── usecase_create.go
│   │   │   ├── usecase_get.go
│   │   │   ├── usecase_list.go
│   │   │   ├── usecase_update.go
│   │   │   └── usecase_delete.go
│   │   └── user/
│   │       ├── dto.go
│   │       ├── usecase_create.go
│   │       ├── usecase_get.go
│   │       ├── usecase_list.go
│   │       ├── usecase_update.go
│   │       ├── usecase_delete.go
│   │       └── usecase_login.go
│   │
│   ├── adapters/                      # Adapters Layer
│   │   ├── http/                      # Driving Adapter (HTTP)
│   │   │   ├── article_handler.go
│   │   │   ├── user_handler.go
│   │   │   ├── router.go
│   │   │   └── middleware.go
│   │   ├── db/                        # Driven Adapter (Database)
│   │   │   ├── article_mysql_repo.go
│   │   │   └── user_mysql_repo.go
│   │   ├── cache/                     # Driven Adapter (Cache)
│   │   │   └── article_redis_cache.go
│   │   └── external/                  # Driven Adapter (External Services)
│   │       └── email_sender.go
│   │
│   └── infrastructure/                # Infrastructure Layer
│       ├── config/
│       │   └── config.go              # Konfigurasi aplikasi
│       ├── database/
│       │   ├── mysql.go               # Koneksi MySQL
│       │   └── redis.go               # Koneksi Redis
│       ├── di/                        # Dependency Injection
│       │   ├── container.go
│       │   ├── article_container.go
│       │   └── user_container.go
│       └── logger/
│           └── logger.go
├── go.mod
├── go.sum
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

**Contoh: `internal/application/article/usecase_get.go`**
```go
// GetArticleUseCase mengorkestrasikan logika untuk mendapatkan artikel
type GetArticleUseCase struct {
    articleRepo domainarticle.Repository  // Port dari domain
    cache       ArticleSingleCache        // Port untuk cache
}

func (uc *GetArticleUseCase) Execute(ctx context.Context, id int64) (*ArticleResponse, error) {
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

**Contoh: `internal/adapters/http/article_handler.go`**
```go
// ArticleHandler adalah driving adapter yang menerima HTTP request
type ArticleHandler struct {
    getUseCase *article.GetArticleUseCase  // Menggunakan use case
}

func (h *ArticleHandler) Get(c *gin.Context) {
    // 1. Parse request
    // 2. Panggil use case
    // 3. Return HTTP response
}
```

#### Driven Adapters (Output)
- **Database Repository**: Implementasi repository untuk MySQL
- **Cache**: Implementasi cache untuk Redis
- **External Services**: Integrasi dengan API eksternal

**Contoh: `internal/adapters/db/article_mysql_repo.go`**
```go
// ArticleMySQLRepository adalah driven adapter yang mengimplementasikan
// domain article.Repository interface
type ArticleMySQLRepository struct {
    db *sql.DB
}

func (r *ArticleMySQLRepository) GetByID(ctx context.Context, id int64) (*article.Article, error) {
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

**Contoh: `internal/infrastructure/di/article_container.go`**
```go
// NewArticleContainer melakukan dependency injection
func NewArticleContainer(database *sql.DB, redisClient *redis.Client) *ArticleContainer {
    // 1. Buat repository (driven adapter)
    articleRepo := db.NewArticleMySQLRepository(database)
    
    // 2. Buat cache (driven adapter)
    articleCache := cache.NewArticleRedisCache(redisClient, 5*time.Minute)
    
    // 3. Buat domain service
    articleService := domainarticle.NewService(articleRepo)
    
    // 4. Buat use cases
    getArticleUseCase := apparticle.NewGetArticleUseCase(articleRepo, articleCache)
    
    // 5. Buat handler (driving adapter)
    articleHandler := http.NewArticleHandler(getArticleUseCase, ...)
    
    return &ArticleContainer{...}
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
```

4. **Setup Database**

Buat database MySQL:
```sql
CREATE DATABASE hexa_go CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

Tabel akan dibuat otomatis saat aplikasi pertama kali dijalankan.

## 🚀 Menjalankan Aplikasi

1. **Jalankan aplikasi**
```bash
go run cmd/api/main.go
```

2. **Aplikasi akan berjalan di**
```
http://localhost:8080
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

# Login
curl -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'

# Create article (dengan token JWT)
curl -X POST http://localhost:8080/api/v1/articles \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "title": "My First Article",
    "content": "This is the content of my article"
  }'

# Get article
curl -X GET http://localhost:8080/api/v1/articles/1

# List articles
curl -X GET http://localhost:8080/api/v1/articles?limit=10&offset=0
```

## 📡 Struktur API

### User Endpoints

- `POST /api/v1/users/register` - Register user baru
- `POST /api/v1/users/login` - Login user
- `GET /api/v1/users/:id` - Get user by ID
- `GET /api/v1/users` - List users
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user

### Article Endpoints

- `POST /api/v1/articles` - Create article (requires authentication)
- `GET /api/v1/articles/:id` - Get article by ID
- `GET /api/v1/articles` - List articles (with pagination)
- `PUT /api/v1/articles/:id` - Update article (requires authentication)
- `DELETE /api/v1/articles/:id` - Delete article (requires authentication)

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

