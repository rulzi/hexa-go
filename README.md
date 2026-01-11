# Hexa-Go: Implementasi Hexagonal Architecture dengan Go

Proyek ini adalah contoh implementasi **Hexagonal Architecture** (Ports and Adapters) menggunakan bahasa pemrograman Go. Aplikasi ini menyediakan API REST untuk manajemen artikel dan pengguna dengan struktur yang terorganisir, mudah diuji, dan dapat dirawat.

> 📖 **Dokumentasi Lengkap**: Untuk dokumentasi detail tentang implementasi Hexagonal Architecture, diagram arsitektur, dan penjelasan mendalam, silakan kunjungi [Wiki Dokumentasi](https://github.com/rulzi/hexa-go/wiki)

## 🏗️ Apa itu Hexagonal Architecture?

**Hexagonal Architecture** (Ports and Adapters) memisahkan logika bisnis dari infrastruktur eksternal:

- **Domain Layer (Core)**: Logika bisnis murni, bebas dari framework
- **Ports**: Interface yang didefinisikan domain
- **Adapters**: Implementasi konkret (HTTP, Database, Cache, dll)

### Keuntungan

✅ **Independensi Framework** - Logika bisnis tidak terikat framework  
✅ **Testabilitas** - Mudah diuji dengan mock dependencies  
✅ **Fleksibilitas** - Ganti teknologi tanpa mengubah domain logic  
✅ **Maintainability** - Kode terorganisir dan mudah dirawat

### Konsep Dasar

```
Driving Adapters (HTTP, CLI, gRPC)
         ↓
Application Layer (Use Cases)
         ↓
Domain Layer (Entities, Ports)
         ↓
Driven Adapters (Database, Cache, External APIs)
```

## 📁 Struktur Proyek

```
hexa-go/
├── cmd/api/main.go              # Entry point
├── internal/
│   ├── domain/                  # Domain Layer (Core)
│   │   ├── user/                # Entities, Ports, Services
│   │   ├── article/
│   │   └── media/
│   ├── application/             # Application Layer
│   │   ├── user/usecase/        # Use Cases
│   │   ├── article/usecase/
│   │   └── media/usecase/
│   ├── adapters/                # Adapters Layer
│   │   ├── http/                # HTTP Handlers (Driving)
│   │   ├── db/                  # Database (Driven)
│   │   ├── auth/                # JWT, Bcrypt (Driven)
│   │   ├── cache/               # Redis Cache (Driven)
│   │   └── storage/             # File Storage (Driven)
│   └── infrastructure/          # Infrastructure
│       ├── database/            # DB Connections
│       └── di/                  # Dependency Injection
└── migration/                   # SQL Migrations
```

## 🔧 Key Points Arsitektur

### Domain Layer
- ✅ **100% Framework-Free** - Hanya standard library Go
- ✅ **Ports didefinisikan di domain** - Repository, TokenGenerator, PasswordHasher, Cache, dll
- ✅ **Entities dengan business logic** - Validasi dan rules bisnis

### Application Layer
- ✅ **Use Cases** - Satu use case = satu operasi bisnis
- ✅ **Menggunakan ports** - Tidak tahu implementasi konkret
- ✅ **DTOs** - Data Transfer Objects untuk komunikasi

### Adapters Layer
- ✅ **Driving Adapters** - HTTP Handlers (Gin)
- ✅ **Driven Adapters** - MySQL, Redis, JWT, Bcrypt, Storage
- ✅ **Mengimplementasikan ports** - Dapat diganti tanpa mengubah domain

### Dependency Flow
```
Infrastructure (DI) → Adapters → Domain Ports ← Application ← HTTP Handlers
```

## 🛠️ Teknologi

- **Go 1.23+**
- **Gin** - Web framework
- **MySQL** - Database
- **Redis** - Cache
- **JWT** - Authentication

## 🚀 Quick Start

### Dengan Docker (Recommended)

```bash
# Clone repository
git clone <repository-url>
cd hexa-go

# Jalankan aplikasi
docker-compose up -d

# Lihat logs
docker-compose logs -f app

# Aplikasi berjalan di http://localhost:8080
```

### Manual Setup

```bash
# Install dependencies
go mod download

# Setup .env file
cp .env.example .env

# Setup database
mysql -u root -p < migration/user.sql
mysql -u root -p < migration/article.sql
mysql -u root -p < migration/media.sql

# Jalankan aplikasi
go run cmd/api/main.go
```

## 📡 API Endpoints

### User
- `POST /api/v1/users/register` - Register (Public)
- `POST /api/v1/users/login` - Login (Public)
- `GET /api/v1/users` - List users (Protected)
- `GET /api/v1/users/:id` - Get user (Protected)

### Article
- `POST /api/v1/articles` - Create (Protected)
- `GET /api/v1/articles` - List (Protected)
- `GET /api/v1/articles/:id` - Get (Protected)
- `PUT /api/v1/articles/:id` - Update (Protected)
- `DELETE /api/v1/articles/:id` - Delete (Protected)

### Media
- `POST /api/v1/media` - Upload (Protected)
- `GET /api/v1/media` - List (Protected)
- `GET /api/v1/media/:id` - Get (Protected)

## 📦 Response Format

```json
{
  "status": "success" | "error",
  "message": "Pesan hasil operasi",
  "data": {} // Optional
}
```

## 🎯 Prinsip Hexagonal Architecture

1. **Dependency Inversion** - Domain tidak bergantung pada adapters
2. **Interface Segregation** - Ports dengan tanggung jawab spesifik
3. **Single Responsibility** - Satu layer = satu tanggung jawab
4. **Open/Closed Principle** - Mudah menambah adapter baru

## 🏆 Key Features

- ✅ **Domain Independence** - 100% bebas dari framework
- ✅ **Ports & Adapters** - Semua ports didefinisikan di domain
- ✅ **Testability** - Setiap layer dapat di-test independen
- ✅ **Flexibility** - Ganti teknologi tanpa mengubah domain

## 📚 Referensi

- 📖 [Wiki Dokumentasi Lengkap](https://github.com/rulzi/hexa-go/wiki)
- [Hexagonal Architecture by Alistair Cockburn](https://alistair.cockburn.us/hexagonal-architecture/)
- [Clean Architecture by Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

## 📝 License

MIT License

## 👤 Author

**Khoirul Afandi**

- Instagram: [@afandi_](https://instagram.com/afandi_)
- LinkedIn: [Khoirul Afandi](https://www.linkedin.com/in/khoirulafandi/)

---

**Dibuat dengan ❤️ menggunakan Go dan Hexagonal Architecture**
