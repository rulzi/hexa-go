# Hexa-Go: Implementasi Hexagonal Architecture dengan Go

Proyek ini adalah contoh implementasi **Hexagonal Architecture** (Ports and Adapters) menggunakan bahasa pemrograman Go. Aplikasi ini menyediakan API REST untuk manajemen artikel dan pengguna dengan struktur yang terorganisir, mudah diuji, dan dapat dirawat.

## 📋 Daftar Isi

- [Apa itu Hexagonal Architecture?](#-apa-itu-hexagonal-architecture)
- [Struktur Proyek](#-struktur-proyek)
- [Komponen Arsitektur](#-komponen-arsitektur)
- [Alur Data](#-alur-data)
- [Teknologi yang Digunakan](#-teknologi-yang-digunakan)
- [Persyaratan](#-persyaratan)
- [Menjalankan dengan Docker](#-menjalankan-dengan-docker)
- [Instalasi dan Konfigurasi](#%EF%B8%8F-instalasi-dan-konfigurasi)
- [Menjalankan Aplikasi](#-menjalankan-aplikasi)
- [Struktur API](#-struktur-api)
- [Format Response Standar](#-format-response-standar)
- [Prinsip-Prinsip Hexagonal Architecture](#-prinsip-prinsip-hexagonal-architecture-dalam-proyek-ini)
- [Arsitektur Implementasi](#-arsitektur-implementasi)
- [Referensi](#-referensi)
- [License](#-license)

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

### Diagram Arsitektur Detail

Diagram berikut menunjukkan arsitektur lengkap dengan semua ports dan adapters:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           DRIVING ADAPTERS                                   │
│  ┌──────────────────────────────────────────────────────────────────────┐   │
│  │  HTTP Adapters (Gin Framework)                                      │   │
│  │  ├── UserHandler    ├── ArticleHandler    ├── MediaHandler         │   │
│  │  └── AuthMiddleware (TokenValidator)                                │   │
│  └──────────────────────────────────────────────────────────────────────┘   │
└───────────────────────────────┬─────────────────────────────────────────────┘
                                │
                                │ Uses
                                ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                        APPLICATION LAYER                                     │
│  ┌──────────────────────────────────────────────────────────────────────┐  │
│  │  Use Cases (Orchestration)                                           │  │
│  │  ├── User Use Cases                                                  │  │
│  │  │   ├── CreateUserUseCase                                           │  │
│  │  │   │   ├── Repository (port)                                       │  │
│  │  │   │   ├── PasswordHasher (port)                                   │  │
│  │  │   │   └── NotificationService (port)                              │  │
│  │  │   ├── LoginUseCase                                                │  │
│  │  │   │   ├── Repository (port)                                       │  │
│  │  │   │   ├── PasswordHasher (port)                                   │  │
│  │  │   │   └── TokenGenerator (port)                                   │  │
│  │  │   └── ...                                                          │  │
│  │  ├── Article Use Cases                                               │  │
│  │  │   ├── GetArticleUseCase                                           │  │
│  │  │   │   ├── Repository (port)                                       │  │
│  │  │   │   └── Cache (port)                                            │  │
│  │  │   └── ...                                                          │  │
│  │  └── Media Use Cases                                                  │  │
│  └──────────────────────────────────────────────────────────────────────┘  │
└───────────────────────────────┬─────────────────────────────────────────────┘
                                │
                                │ Depends on
                                ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                            DOMAIN LAYER                                       │
│  ┌──────────────────────────────────────────────────────────────────────┐  │
│  │  Entities & Business Logic                                           │  │
│  │  ├── User, Article, Media (entities with validation)                │  │
│  │  └── Domain Services (business logic)                                │  │
│  │                                                                       │  │
│  │  PORTS (Interfaces) - Defined by Domain                              │  │
│  │  ├── Repository Ports                                               │  │
│  │  │   ├── user.Repository                                             │  │
│  │  │   ├── article.Repository                                          │  │
│  │  │   └── media.Repository                                            │  │
│  │  ├── Authentication Ports                                            │  │
│  │  │   ├── TokenGenerator                                              │  │
│  │  │   ├── TokenValidator                                              │  │
│  │  │   └── PasswordHasher                                               │  │
│  │  ├── Infrastructure Ports                                            │  │
│  │  │   ├── article.Cache                                               │  │
│  │  │   ├── media.Storage                                                │  │
│  │  │   └── user.NotificationService                                     │  │
│  └──────────────────────────────────────────────────────────────────────┘  │
└───────────────────────────────┬─────────────────────────────────────────────┘
                                │
                                │ Implemented by
                                ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                          DRIVEN ADAPTERS                                      │
│  ┌──────────────────────────────────────────────────────────────────────┐   │
│  │  Database Adapters                                                   │   │
│  │  ├── MySQLRepository (implements Repository ports)                 │   │
│  │  │   ├── user.Repository → user.Repository                           │   │
│  │  │   ├── article.Repository → article.Repository                    │   │
│  │  │   └── media.Repository → media.Repository                         │   │
│  │                                                                      │   │
│  │  Authentication Adapters                                            │   │
│  │  ├── JWTAdapter → TokenGenerator, TokenValidator                   │   │
│  │  └── BcryptPasswordHasher → PasswordHasher                         │   │
│  │                                                                      │   │
│  │  Cache Adapters                                                     │   │
│  │  ├── RedisCache (DTO-based)                                         │   │
│  │  └── DomainCacheAdapter → article.Cache                             │   │
│  │                                                                      │   │
│  │  External Service Adapters                                          │   │
│  │  ├── EmailSenderImpl → NotificationService                          │   │
│  │  └── LocalStorage → media.Storage                                    │   │
│  └──────────────────────────────────────────────────────────────────────┘   │
└───────────────────────────────┬─────────────────────────────────────────────┘
                                │
                                │ Managed by
                                ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                        INFRASTRUCTURE LAYER                                   │
│  ┌──────────────────────────────────────────────────────────────────────┐  │
│  │  Dependency Injection Containers                                     │  │
│  │  ├── di/container.go (main container)                               │  │
│  │  ├── di/user/container.go                                            │  │
│  │  ├── di/article/container.go                                        │  │
│  │  └── di/media/container.go                                           │  │
│  │                                                                      │  │
│  │  Infrastructure Setup                                                │  │
│  │  ├── database/ (MySQL, Redis connections)                          │  │
│  │  ├── config/ (configuration management)                              │  │
│  │  └── logger/ (logging setup)                                         │  │
│  └──────────────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Diagram Dependency Flow

```
┌─────────────────────────────────────────────────────────────────────────┐
│                        DEPENDENCY FLOW                                    │
│                                                                           │
│  Infrastructure Layer                                                     │
│  ┌───────────────────────────────────────────────────────────────────┐  │
│  │  DI Container wires:                                               │  │
│  │  Adapters → Use Cases → Domain Ports                              │  │
│  └───────────────────────────────────────────────────────────────────┘  │
│                              │                                            │
│                              │ Creates & Wires                            │
│                              ▼                                            │
│  ┌───────────────────────────────────────────────────────────────────┐  │
│  │  HTTP Handler (Driving Adapter)                                    │  │
│  │  ┌─────────────────────────────────────────────────────────────┐  │  │
│  │  │  Handler.Create() → UseCase.Execute()                       │  │  │
│  │  └─────────────────────────────────────────────────────────────┘  │  │
│  └───────────────────────────────────────────────────────────────────┘  │
│                              │                                            │
│                              │ Calls                                      │
│                              ▼                                            │
│  ┌───────────────────────────────────────────────────────────────────┐  │
│  │  Use Case (Application Layer)                                     │  │
│  │  ┌─────────────────────────────────────────────────────────────┐  │  │
│  │  │  CreateUserUseCase.Execute()                                 │  │  │
│  │  │  ├── passwordHasher.Hash() [port]                           │  │  │
│  │  │  ├── repository.Create() [port]                             │  │  │
│  │  │  └── notificationService.SendWelcomeEmail() [port]         │  │  │
│  │  └─────────────────────────────────────────────────────────────┘  │  │
│  └───────────────────────────────────────────────────────────────────┘  │
│                              │                                            │
│                              │ Uses Ports                                 │
│                              ▼                                            │
│  ┌───────────────────────────────────────────────────────────────────┐  │
│  │  Domain Ports (Interfaces)                                        │  │
│  │  ┌─────────────────────────────────────────────────────────────┐  │  │
│  │  │  PasswordHasher, Repository, NotificationService          │  │  │
│  │  └─────────────────────────────────────────────────────────────┘  │  │
│  └───────────────────────────────────────────────────────────────────┘  │
│                              ▲                                            │
│                              │ Implemented by                             │
│                              │                                            │
│  ┌───────────────────────────────────────────────────────────────────┐  │
│  │  Adapters (Driven)                                                │  │
│  │  ┌─────────────────────────────────────────────────────────────┐  │  │
│  │  │  BcryptPasswordHasher → PasswordHasher                      │  │  │
│  │  │  MySQLRepository → Repository                               │  │  │
│  │  │  EmailSenderImpl → NotificationService                      │  │  │
│  │  └─────────────────────────────────────────────────────────────┘  │  │
│  └───────────────────────────────────────────────────────────────────┘  │
│                                                                           │
│  ✅ Dependency Direction: Adapters → Application → Domain                │
│  ✅ Domain is independent (no dependencies on other layers)              │
└─────────────────────────────────────────────────────────────────────────┘
```

### Diagram User Domain (Detail)

Contoh detail implementasi untuk User domain:

```
┌─────────────────────────────────────────────────────────────────────────┐
│                          USER DOMAIN FLOW                                │
│                                                                           │
│  ┌───────────────────────────────────────────────────────────────────┐  │
│  │  HTTP Handler (Driving Adapter)                                    │  │
│  │  ┌─────────────────────────────────────────────────────────────┐  │  │
│  │  │  POST /users/register                                       │  │  │
│  │  │  → CreateUserUseCase.Execute()                              │  │  │
│  │  └─────────────────────────────────────────────────────────────┘  │  │
│  └───────────────────────────────────────────────────────────────────┘  │
│                              │                                            │
│                              ▼                                            │
│  ┌───────────────────────────────────────────────────────────────────┐  │
│  │  CreateUserUseCase (Application)                                  │  │
│  │  ┌─────────────────────────────────────────────────────────────┐  │  │
│  │  │  Dependencies (Ports):                                      │  │  │
│  │  │  ├── Repository (port)                                      │  │  │
│  │  │  ├── PasswordHasher (port)                                  │  │  │
│  │  │  └── NotificationService (port)                             │  │  │
│  │  │                                                             │  │  │
│  │  │  Flow:                                                      │  │  │
│  │  │  1. Check email exists (Repository)                         │  │  │
│  │  │  2. Hash password (PasswordHasher)                           │  │  │
│  │  │  3. Create user (Repository)                                 │  │  │
│  │  │  4. Send welcome email (NotificationService)                 │  │  │
│  │  └─────────────────────────────────────────────────────────────┘  │  │
│  └───────────────────────────────────────────────────────────────────┘  │
│                              │                                            │
│                              │ Uses Ports                                 │
│                              ▼                                            │
│  ┌───────────────────────────────────────────────────────────────────┐  │
│  │  Domain Ports (Defined in domain/user/)                          │  │
│  │  ┌─────────────────────────────────────────────────────────────┐  │  │
│  │  │  • Repository interface                                     │  │  │
│  │  │  • PasswordHasher interface                                 │  │  │
│  │  │  • NotificationService interface                            │  │  │
│  │  │  • TokenGenerator interface                                 │  │  │
│  │  │  • TokenValidator interface                                 │  │  │
│  │  └─────────────────────────────────────────────────────────────┘  │  │
│  └───────────────────────────────────────────────────────────────────┘  │
│                              ▲                                            │
│                              │ Implemented by                             │
│                              │                                            │
│  ┌───────────────────────────────────────────────────────────────────┐  │
│  │  Adapters (Implementations)                                       │  │
│  │  ┌─────────────────────────────────────────────────────────────┐  │  │
│  │  │  • Repository → Repository                                      │  │  │
│  │  │  • BcryptPasswordHasher → PasswordHasher                    │  │  │
│  │  │  • EmailSenderImpl → NotificationService                     │  │  │
│  │  │  • JWTAdapter → TokenGenerator, TokenValidator              │  │  │
│  │  └─────────────────────────────────────────────────────────────┘  │  │
│  └───────────────────────────────────────────────────────────────────┘  │
│                              │                                            │
│                              │ Connects to                                 │
│                              ▼                                            │
│  ┌───────────────────────────────────────────────────────────────────┐  │
│  │  External Systems                                                 │  │
│  │  ┌─────────────────────────────────────────────────────────────┐  │  │
│  │  │  • MySQL Database                                            │  │  │
│  │  │  • Email Service (SMTP/SendGrid/etc)                         │  │  │
│  │  └─────────────────────────────────────────────────────────────┘  │  │
│  └───────────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────┘
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
│   │   │   ├── cache.go               # Port (Interface) untuk cache
│   │   │   └── errors.go              # Domain errors
│   │   ├── user/
│   │   │   ├── entity.go              # Entity User
│   │   │   ├── repository.go          # Port (Interface) untuk repository
│   │   │   ├── service.go             # Domain service
│   │   │   ├── token.go               # Ports: TokenGenerator, TokenValidator
│   │   │   ├── password.go            # Port: PasswordHasher
│   │   │   ├── notification.go        # Port: NotificationService
│   │   │   └── errors.go              # Domain errors
│   │   └── media/
│   │       ├── entity.go
│   │       ├── repository.go
│   │       ├── service.go
│   │       ├── storage.go             # Port (Interface) untuk storage
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
│   │   │   ├── middleware/            # HTTP Middlewares
│   │   │   │   ├── auth.go            # JWT Authentication middleware
│   │   │   │   ├── cors.go            # CORS middleware
│   │   │   │   ├── recovery.go        # Panic recovery middleware
│   │   │   │   └── setup.go           # Middleware setup
│   │   │   ├── response/              # Standard response format
│   │   │   │   └── response.go
│   │   │   └── router.go
│   │   ├── db/                        # Driven Adapter (Database)
│   │   │   ├── article/
│   │   │   │   └── repository.go
│   │   │   ├── user/
│   │   │   │   └── repository.go
│   │   │   └── media/
│   │   │       └── repository.go
│   │   ├── auth/                      # Driven Adapter (Authentication)
│   │   │   ├── jwt_adapter.go         # Implements TokenGenerator, TokenValidator
│   │   │   └── bcrypt_adapter.go      # Implements PasswordHasher
│   │   ├── cache/                     # Driven Adapter (Cache)
│   │   │   └── article/
│   │   │       ├── redis_cache.go     # Redis cache implementation
│   │   │       └── domain_cache_adapter.go  # Adapter untuk domain cache port
│   │   ├── storage/                   # Driven Adapter (File Storage)
│   │   │   └── media/
│   │   │       └── local_storage.go
│   │   └── external/                  # Driven Adapter (External Services)
│   │       └── user/
│   │           └── email_sender.go    # Implements NotificationService
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

**Ports yang Didefinisikan di Domain:**

1. **Repository Ports** (Persistence):
   - `user.Repository` - User persistence
   - `article.Repository` - Article persistence
   - `media.Repository` - Media persistence

2. **Authentication Ports**:
   - `user.TokenGenerator` - Generate authentication tokens
   - `user.TokenValidator` - Validate authentication tokens
   - `user.PasswordHasher` - Hash and verify passwords

3. **Infrastructure Ports**:
   - `article.Cache` - Article caching
   - `media.Storage` - File storage
   - `user.NotificationService` - Send notifications (emails)

**Contoh: `internal/domain/user/token.go`**
```go
// TokenGenerator adalah port (interface) yang didefinisikan oleh domain
// Domain tidak peduli apakah implementasinya JWT, OAuth, atau lainnya
type TokenGenerator interface {
    Generate(userID int64, email string) (string, error)
}

type TokenValidator interface {
    Validate(token string) (*TokenClaims, error)
}
```

**Contoh: `internal/domain/user/password.go`**
```go
// PasswordHasher adalah port untuk password operations
// Domain tidak peduli apakah implementasinya bcrypt, argon2, atau lainnya
type PasswordHasher interface {
    Hash(password string) (string, error)
    Verify(hashedPassword, password string) bool
}
```

**Karakteristik:**
- ✅ **100% bebas dari framework atau teknologi eksternal**
- ✅ Hanya bergantung pada standard library Go (`context`, `time`, `errors`)
- ✅ Hanya berisi logika bisnis murni
- ✅ Mendefinisikan kontrak (ports) yang harus dipenuhi oleh adapters
- ✅ Domain services menggunakan ports, bukan concrete implementations

### 2. Application Layer (`internal/application/`)

**Application Layer** berisi use cases yang mengorkestrasikan domain logic:

- **Use Cases**: Setiap use case mewakili satu operasi bisnis
- **DTOs**: Data Transfer Objects untuk komunikasi antar layer

**Contoh: `internal/application/user/usecase/create.go`**
```go
// CreateUserUseCase mengorkestrasikan logika untuk membuat user baru
type CreateUserUseCase struct {
    userRepo            domainuser.Repository          // Port dari domain
    passwordHasher      domainuser.PasswordHasher      // Port dari domain
    notificationService domainuser.NotificationService // Port dari domain
}

func (uc *CreateUserUseCase) Execute(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
    // 1. Check email exists (Repository port)
    // 2. Hash password (PasswordHasher port)
    // 3. Create user (Repository port)
    // 4. Send welcome email (NotificationService port)
    // 5. Return response DTO
}
```

**Contoh: `internal/application/article/usecase/get.go`**
```go
// GetArticleUseCase mengorkestrasikan logika untuk mendapatkan artikel
type GetArticleUseCase struct {
    articleRepo domainarticle.Repository  // Port dari domain
    cache       domainarticle.Cache        // Port dari domain
}

func (uc *GetArticleUseCase) Execute(ctx context.Context, id int64) (*dto.ArticleResponse, error) {
    // 1. Cek cache dulu (Cache port)
    // 2. Jika tidak ada, ambil dari repository (Repository port)
    // 3. Simpan ke cache (Cache port)
    // 4. Return response DTO
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

**Contoh: `internal/adapters/db/article/repository.go`**
```go
// MySQLRepository adalah driven adapter yang mengimplementasikan
// domain article.Repository interface
type MySQLRepository struct {
    db *sql.DB
}

func (r *MySQLRepository) GetByID(ctx context.Context, id int64) (*article.Article, error) {
    // Implementasi konkret menggunakan MySQL
    // Mengkonversi infrastructure errors ke domain errors
    if err == sql.ErrNoRows {
        return nil, domainarticle.ErrArticleNotFound
    }
}
```

**Contoh: `internal/adapters/auth/jwt_adapter.go`**
```go
// JWTAdapter adalah driven adapter yang mengimplementasikan
// domain user.TokenGenerator dan user.TokenValidator interfaces
type JWTAdapter struct {
    secret     string
    expiration int
}

func (a *JWTAdapter) Generate(userID int64, email string) (string, error) {
    // Implementasi konkret menggunakan JWT library
    // Domain tidak tahu bahwa ini menggunakan JWT
}

func (a *JWTAdapter) Validate(tokenString string) (*domainuser.TokenClaims, error) {
    // Implementasi konkret menggunakan JWT library
}
```

**Contoh: `internal/adapters/auth/bcrypt_adapter.go`**
```go
// BcryptPasswordHasher adalah driven adapter yang mengimplementasikan
// domain user.PasswordHasher interface
type BcryptPasswordHasher struct{}

func (h *BcryptPasswordHasher) Hash(password string) (string, error) {
    // Implementasi konkret menggunakan bcrypt library
    // Domain tidak tahu bahwa ini menggunakan bcrypt
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

## 🐳 Menjalankan dengan Docker

Cara termudah untuk menjalankan aplikasi adalah menggunakan Docker Compose, yang akan menjalankan semua dependencies (MySQL, Redis) dan aplikasi Go secara otomatis.

### Persyaratan Docker

- Docker 20.10 atau lebih tinggi
- Docker Compose 2.0 atau lebih tinggi

### 🚀 Quick Start dengan Docker

1. **Clone repository**
```bash
git clone <repository-url>
cd hexa-go
```

2. **Jalankan aplikasi dengan Docker Compose**
```bash
docker-compose up -d
```

3. **Tunggu beberapa saat hingga semua service berjalan**
```bash
# Lihat status containers
docker-compose ps

# Lihat logs
docker-compose logs -f app
```

4. **Aplikasi akan berjalan di**
```
http://localhost:8080
```

### 📦 Services yang Berjalan

Docker Compose akan menjalankan 3 services:

- **hexa-app**: Aplikasi Go (port 8080)
- **hexa-mysql**: MySQL 8.0 database (port 3306)
- **hexa-redis**: Redis cache (port 6379)

### 🔧 Environment Variables untuk Docker

Environment variables sudah dikonfigurasi dalam `docker-compose.yml`. Untuk kustomisasi, Anda bisa:

1. **Edit docker-compose.yml** (langsung ubah environment section)
2. **Buat file `.env`** di root directory
3. **Gunakan environment file custom**
```bash
docker-compose --env-file .env.custom up -d
```

### 📋 Docker Commands yang Berguna

```bash
# Jalankan aplikasi
docker-compose up -d

# Lihat logs aplikasi
docker-compose logs -f app

# Lihat logs database
docker-compose logs -f mysql

# Lihat logs Redis
docker-compose logs -f redis

# Restart aplikasi
docker-compose restart app

# Stop semua services
docker-compose down

# Stop dan hapus volumes (menghapus data)
docker-compose down -v

# Rebuild aplikasi (jika ada perubahan kode)
docker-compose build app
docker-compose up -d app

# Eksekusi command di dalam container
docker-compose exec app sh
docker-compose exec app go run cmd/api/main.go

# Lihat status services
docker-compose ps

# Lihat resource usage
docker-compose top
```

### 🔍 Health Checks

Aplikasi memiliki health check endpoint yang bisa dipantau:

```bash
# Health check aplikasi
curl http://localhost:8080/health

# Health check database via docker-compose
docker-compose exec mysql mysqladmin ping -h localhost
```

### 🗄️ Volume Management

Data akan tersimpan dalam named volumes:

- **mysql_data**: Data database MySQL
- **redis_data**: Data Redis
- **storage_data**: File upload storage

```bash
# Lihat volumes
docker volume ls | grep hexa

# Backup database
docker-compose exec mysql mysqldump -u hexa_user -phexapassword123 hexa_go > backup.sql

# Restore database
docker-compose exec -T mysql mysql -u hexa_user -phexapassword123 hexa_go < backup.sql
```

### 🌐 Akses Database dari Host

Untuk development, Anda bisa akses MySQL dari host:

```bash
# Connect ke MySQL dari host
mysql -h 127.0.0.1 -P 3306 -u hexa_user -phexapassword123 hexa_go

# Connect ke Redis dari host
redis-cli -h 127.0.0.1 -p 6379
```

### 🛠️ Development dengan Docker

#### Development Mode (dengan code reload)

Untuk development, Anda bisa mount source code dan gunakan hot reload:

1. **Edit docker-compose.yml** untuk development:

```yaml
# Tambahkan volume untuk development
app:
  build:
    context: .
    dockerfile: Dockerfile
  volumes:
    - .:/app  # Mount source code
    - storage_data:/app/storage
  environment:
    # ... environment variables
```

2. **Gunakan air untuk hot reload**:

```bash
# Install air
go install github.com/cosmtrek/air@latest

# Run dengan air
air
```

#### Production Mode

Untuk production, gunakan multi-stage build yang sudah dikonfigurasi:

```bash
# Build production image
docker build -t hexa-go:latest .

# Jalankan container
docker run -d \
  --name hexa-go-prod \
  -p 8080:8080 \
  -e DB_HOST=your-db-host \
  -e DB_USER=your-db-user \
  -e DB_PASSWORD=your-db-password \
  -e DB_NAME=your-db-name \
  -e REDIS_HOST=your-redis-host \
  -e JWT_SECRET=your-production-secret \
  hexa-go:latest
```

### 🐛 Troubleshooting Docker

**Container tidak bisa start:**
```bash
# Lihat logs detail
docker-compose logs app
docker-compose logs mysql
docker-compose logs redis
```

**Database connection error:**
```bash
# Pastikan MySQL sudah healthy
docker-compose exec mysql mysqladmin ping -h localhost
docker-compose logs mysql
```

**Port sudah digunakan:**
```bash
# Cek port yang digunakan
lsof -i :8080
lsof -i :3306
lsof -i :6379

# Stop container yang menggunakan port
docker-compose down
```

**Reset semua data:**
```bash
# Hapus semua containers dan volumes
docker-compose down -v
docker system prune -f

# Jalankan ulang
docker-compose up -d
```

### 🔒 Production Security Notes

1. **Ganti default passwords** dalam `docker-compose.yml`
2. **Gunakan environment files** untuk secrets
3. **Gunakan Docker secrets** untuk data sensitif
4. **Setup SSL/TLS** untuk production
5. **Gunakan specific image tags** bukan `latest`

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

5. **Test API**

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

## 🏆 Arsitektur Implementasi

### Domain Independence (100% Framework-Free)

Domain layer **100% bebas dari framework dan library eksternal**:

- ✅ **Tidak ada import** `gin`, `sql`, `redis`, `jwt`, `bcrypt` di domain layer
- ✅ **Hanya standard library**: `context`, `time`, `errors`
- ✅ **Semua ports didefinisikan di domain**: TokenGenerator, PasswordHasher, NotificationService, Cache
- ✅ **Domain services menggunakan ports**: Tidak ada concrete implementations

### Ports & Adapters Pattern

**Semua ports didefinisikan di domain layer:**

| Port | Lokasi | Implementasi |
|------|--------|--------------|
| `user.Repository` | `domain/user/repository.go` | `adapters/db/user/repository.go` |
| `user.TokenGenerator` | `domain/user/token.go` | `adapters/auth/jwt_adapter.go` |
| `user.TokenValidator` | `domain/user/token.go` | `adapters/auth/jwt_adapter.go` |
| `user.PasswordHasher` | `domain/user/password.go` | `adapters/auth/bcrypt_adapter.go` |
| `user.NotificationService` | `domain/user/notification.go` | `adapters/external/user/email_sender.go` |
| `article.Repository` | `domain/article/repository.go` | `adapters/db/article/repository.go` |
| `article.Cache` | `domain/article/cache.go` | `adapters/cache/article/domain_cache_adapter.go` |
| `media.Storage` | `domain/media/storage.go` | `adapters/storage/media/local_storage.go` |

### Dependency Flow

```
Infrastructure (DI Container)
    ↓ wires
Adapters (Concrete Implementations)
    ↓ implements
Domain Ports (Interfaces)
    ↑ used by
Application Layer (Use Cases)
    ↑ called by
Driving Adapters (HTTP Handlers)
```

**Key Points:**
- ✅ Domain tidak bergantung pada layer lain
- ✅ Application hanya bergantung pada domain ports
- ✅ Adapters mengimplementasikan domain ports
- ✅ DI container meng-wire semua dependencies

### Testability

Dengan arsitektur ini, setiap layer dapat di-test secara independen:

- **Domain**: Test entities dan business logic tanpa dependencies
- **Application**: Mock domain ports untuk test use cases
- **Adapters**: Test implementasi ports secara terpisah
- **Integration**: Test end-to-end dengan real adapters

**Contoh Test Domain:**
```go
// Domain dapat di-test tanpa framework
func TestUser_Validate(t *testing.T) {
    user := &domainuser.User{
        Name:     "John",
        Email:    "john@example.com",
        Password: "hashed",
    }
    err := user.Validate()
    assert.NoError(t, err)
}
```

**Contoh Test Use Case dengan Mock:**
```go
// Use case dapat di-test dengan mock ports
func TestCreateUserUseCase(t *testing.T) {
    mockRepo := &MockRepository{}
    mockPasswordHasher := &MockPasswordHasher{}
    mockNotification := &MockNotificationService{}
    
    useCase := usecase.NewCreateUserUseCase(
        mockRepo,
        mockPasswordHasher,
        mockNotification,
    )
    // Test use case...
}
```

## 📚 Referensi

- [Hexagonal Architecture by Alistair Cockburn](https://alistair.cockburn.us/hexagonal-architecture/)
- [Clean Architecture by Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

## 📝 License

MIT License

---

**Dibuat dengan ❤️ menggunakan Go dan Hexagonal Architecture**
