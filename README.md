# Hexa-Go

[![Go Version](https://img.shields.io/badge/Go-1.26.4-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![CI](https://github.com/rulzi/hexa-go/actions/workflows/ci.yml/badge.svg)](https://github.com/rulzi/hexa-go/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**Hexa-Go** adalah proyek referensi REST API berbasis **Hexagonal Architecture** (Ports & Adapters) di Go. Proyek ini ditujukan sebagai **boilerplate pembelajaran dan starting point** — bukan framework siap produksi — dengan tiga domain bisnis: **User**, **Article**, dan **Media**. Domain layer bebas dari framework; HTTP, MySQL, Redis, JWT, dan storage dihubungkan lewat adapter yang dapat diganti.

> Dokumentasi arsitektur lebih mendalam tersedia di [Wiki](https://github.com/rulzi/hexa-go/wiki).

---

## Arsitektur

Alur dependency mengikuti prinsip **dependency inversion**: domain mendefinisikan port (interface), adapter mengimplementasikannya, dan use case hanya bergantung pada port.

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         Driving Adapters (Inbound)                      │
│  HTTP Handlers (user/article/media/health) + Middleware (auth, CORS,   │
│  logging, recovery, request-id)                                         │
└───────────────────────────────────┬─────────────────────────────────────┘
                                    │ calls
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                         Application Layer                               │
│  Use Cases (user/article/media) + DTOs                                  │
└───────────────────────────────────┬─────────────────────────────────────┘
                                    │ uses ports
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                         Domain Layer (Core)                               │
│  Entities · Domain Services · Ports (Repository, Cache, Storage,       │
│  TokenGenerator, PasswordHasher, NotificationService)                     │
└───────────────────────────────────┬─────────────────────────────────────┘
                                    │ implemented by
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                         Driven Adapters (Outbound)                      │
│  MySQL Repositories · Redis Cache (article list) · JWT/Bcrypt Auth       │
│  Local/S3 Storage · Email Notification (stub)                           │
└─────────────────────────────────────────────────────────────────────────┘

Wiring: cmd/api/main.go → infrastructure/di → per-domain containers
        (user · article · media) → http.Router
```

### Wiring per domain

| Domain  | Repository        | Cache              | Storage        | Auth / Lainnya              |
|---------|-------------------|--------------------|----------------|-----------------------------|
| User    | MySQL             | —                  | —              | JWT, Bcrypt, Email sender   |
| Article | MySQL             | Redis (opsional)   | —              | —                           |
| Media   | MySQL             | —                  | Local / S3     | —                           |

Redis bersifat opsional. Jika koneksi gagal, aplikasi tetap berjalan tanpa cache artikel.

---

## Struktur Folder

```
hexa-go/
├── cmd/api/
│   └── main.go                         # Entry point: config, DB, Redis, DI, HTTP server
├── internal/
│   ├── domain/                         # Core — tanpa dependency framework
│   │   ├── user/                       # Entity, port, domain service
│   │   ├── article/
│   │   ├── media/
│   │   └── errs/                       # Domain error types & mapping helpers
│   ├── application/                    # Use cases & DTO
│   │   ├── user/usecase/
│   │   ├── article/usecase/
│   │   └── media/usecase/
│   ├── adapters/
│   │   ├── http/                       # Driving: handlers, router, middleware, response
│   │   │   ├── user/ article/ media/ health/
│   │   │   ├── middleware/             # auth, cors, logging, recovery, request-id
│   │   │   ├── errmapper/              # Domain error → HTTP status
│   │   │   └── response/               # Standard JSON response envelope
│   │   ├── repository/                 # Driven: MySQL implementations
│   │   ├── cache/article/              # Driven: Redis cache + domain adapter
│   │   ├── auth/                       # Driven: JWT & Bcrypt
│   │   ├── storage/media/              # Driven: local filesystem & AWS S3
│   │   ├── external/user/              # Driven: email notification (stub)
│   │   ├── contextkey/                 # Request-scoped context keys
│   │   └── logging/                    # Repository query logging helper
│   └── infrastructure/
│       ├── config/                     # Env loading & validation
│       ├── database/                   # MySQL & Redis connection factories
│       ├── di/                         # Root DI + per-domain containers
│       └── logger/                     # Logrus-based structured logger
├── migration/                          # SQL schema (auto-run via Docker init)
│   ├── 001_user.sql
│   ├── 002_article.sql
│   └── 003_media.sql
├── .github/workflows/ci.yml            # Lint, test, coverage, Docker publish
├── docker-compose.yml                  # MySQL 8, Redis 7, app
├── Dockerfile
├── Makefile
├── .env.example                        # Local development
└── .env.docker.example                 # Docker Compose
```

---

## Tech Stack

Library berikut tercatat di `go.mod` dan digunakan di kode produksi:

| Library | Versi | Peran |
|---------|-------|-------|
| [gin-gonic/gin](https://github.com/gin-gonic/gin) | v1.12.0 | HTTP router & middleware |
| [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql) | v1.7.1 | Driver MySQL |
| [redis/go-redis/v9](https://github.com/redis/go-redis) | v9.17.2 | Client Redis (cache artikel) |
| [golang-jwt/jwt/v5](https://github.com/golang-jwt/jwt) | v5.3.1 | Generate & validasi JWT |
| [golang.org/x/crypto](https://pkg.go.dev/golang.org/x/crypto) | v0.53.0 | Bcrypt password hashing |
| [aws/aws-sdk-go-v2](https://github.com/aws/aws-sdk-go-v2) (S3) | v1.42.0 / v1.103.3 | Storage S3 (opsional) |
| [sirupsen/logrus](https://github.com/sirupsen/logrus) | v1.9.3 | Structured logging |
| [joho/godotenv](https://github.com/joho/godotenv) | v1.5.1 | Load `.env` lokal |

**Testing & dev tools:**

| Library | Versi | Peran |
|---------|-------|-------|
| [stretchr/testify](https://github.com/stretchr/testify) | v1.11.1 | Assertions & mocks |
| [DATA-DOG/go-sqlmock](https://github.com/DATA-DOG/go-sqlmock) | v1.5.2 | Mock database |
| [alicebob/miniredis/v2](https://github.com/alicebob/miniredis) | v2.35.0 | In-memory Redis untuk test |

---

## Prerequisites

| Kebutuhan | Versi |
|-----------|-------|
| Go | **1.26.4** (sesuai `go.mod`) |
| Docker & Docker Compose | Docker 20+ direkomendasikan |
| MySQL | 8.0 (via Docker atau lokal) |
| Redis | 7.x (opsional, via Docker atau lokal) |

### Environment Variables

**Lokal** — salin dari `.env.example`:

```bash
cp .env.example .env
```

| Variable | Wajib | Keterangan |
|----------|-------|------------|
| `DEBUG` | — | `true` = dev mode; `DB_PASSWORD` & `JWT_SECRET` opsional |
| `SERVER_HOST`, `SERVER_PORT` | — | Default: `0.0.0.0:8080` |
| `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` | `DB_PASSWORD` jika `DEBUG=false` | Koneksi MySQL |
| `REDIS_HOST`, `REDIS_PORT`, `REDIS_PASSWORD` | — | Redis opsional |
| `JWT_SECRET`, `JWT_EXPIRATION` | `JWT_SECRET` (min 32 char) jika `DEBUG=false` | Autentikasi JWT |
| `STORAGE_DRIVER` | — | `local` (default) atau `s3` |
| `STORAGE_BASE_PATH`, `STORAGE_BASE_URL` | — | Path & URL publik file lokal |
| `S3_BUCKET`, `S3_REGION`, `S3_ENDPOINT`, `S3_USE_PATH_STYLE` | Jika `STORAGE_DRIVER=s3` | Konfigurasi S3 |

**Docker** — salin dari `.env.docker.example`:

```bash
cp .env.docker.example .env.docker
```

Ubah `MYSQL_ROOT_PASSWORD`, `MYSQL_PASSWORD`, `DB_PASSWORD`, dan `JWT_SECRET` sebelum menjalankan. Generate secret aman:

```bash
openssl rand -base64 32
```

---

## Getting Started

### Opsi 1: Docker (direkomendasikan)

```bash
git clone https://github.com/rulzi/hexa-go.git
cd hexa-go

# Siapkan environment Docker
cp .env.docker.example .env.docker
# Edit .env.docker — set password & JWT_SECRET

# Jalankan stack (MySQL + Redis + App)
docker compose up -d

# Cek status
docker compose ps
docker compose logs -f app
```

- API: `http://localhost:8080`
- Health check: `GET http://localhost:8080/health`
- Migrasi SQL di `./migration/` dijalankan otomatis saat container MySQL pertama kali start

Atau via Makefile:

```bash
make docker-up      # docker compose up -d
make docker-logs    # tail logs
make docker-down    # stop containers
```

### Opsi 2: Manual (lokal)

```bash
git clone https://github.com/rulzi/hexa-go.git
cd hexa-go

go mod download

# Environment
cp .env.example .env
# Pastikan MySQL & Redis berjalan, sesuaikan DB_* dan REDIS_*

# Migrasi database
mysql -u root -p hexa_go < migration/001_user.sql
mysql -u root -p hexa_go < migration/002_article.sql
mysql -u root -p hexa_go < migration/003_media.sql

# Jalankan
go run cmd/api/main.go
# atau: make run
```

Dengan `DEBUG=true`, aplikasi bisa langsung jalan tanpa mengisi `DB_PASSWORD` dan `JWT_SECRET` (menggunakan default dev).

---

## API Endpoints

Base URL: `http://localhost:8080`

### Autentikasi

Endpoint **Protected** memerlukan header:

```
Authorization: Bearer <jwt_token>
```

Token didapat dari `POST /api/v1/users/login` atau `POST /api/v1/users/register`.

### Health

| Method | Path | Auth | Deskripsi |
|--------|------|------|-----------|
| `GET` | `/health` | Public | Cek koneksi MySQL & Redis |

### User

| Method | Path | Auth | Deskripsi |
|--------|------|------|-----------|
| `POST` | `/api/v1/users/register` | Public | Registrasi pengguna baru |
| `POST` | `/api/v1/users/login` | Public | Login, mengembalikan JWT |
| `POST` | `/api/v1/users` | Protected | Buat pengguna (admin) |
| `GET` | `/api/v1/users` | Protected | Daftar pengguna (`?limit=&offset=`) |
| `GET` | `/api/v1/users/:id` | Protected | Detail pengguna |
| `PUT` | `/api/v1/users/:id` | Protected | Update pengguna |
| `DELETE` | `/api/v1/users/:id` | Protected | Hapus pengguna |

### Article

| Method | Path | Auth | Deskripsi |
|--------|------|------|-----------|
| `POST` | `/api/v1/articles` | Protected | Buat artikel |
| `GET` | `/api/v1/articles` | Protected | Daftar artikel (`?limit=&offset=`) |
| `GET` | `/api/v1/articles/:id` | Protected | Detail artikel |
| `PUT` | `/api/v1/articles/:id` | Protected | Update artikel |
| `DELETE` | `/api/v1/articles/:id` | Protected | Hapus artikel |

### Media

| Method | Path | Auth | Deskripsi |
|--------|------|------|-----------|
| `POST` | `/api/v1/media` | Protected | Upload media (multipart) |
| `GET` | `/api/v1/media` | Protected | Daftar media (`?limit=&offset=`) |
| `GET` | `/api/v1/media/:id` | Protected | Detail media |
| `PUT` | `/api/v1/media/:id` | Protected | Update metadata media |
| `DELETE` | `/api/v1/media/:id` | Protected | Hapus media |
| `GET` | `/api/v1/media/files/*` | Public | Akses file statis (local storage) |

### Response Format

Semua endpoint API mengembalikan envelope JSON:

```json
{
  "status": "success",
  "message": "Pesan hasil operasi",
  "data": {}
}
```

Error:

```json
{
  "status": "error",
  "message": "Deskripsi error"
}
```

---

## Testing

```bash
# Unit test semua paket (dengan race detector)
go test ./... -race -v

# Dengan laporan coverage
go test ./... -race -coverprofile=coverage.out -covermode=atomic
go tool cover -func=coverage.out

# Lihat coverage per paket (HTML)
go tool cover -html=coverage.out
```

Via Makefile:

```bash
make test    # menjalankan lint + test dengan coverage.out
```

### Coverage

CI ([`.github/workflows/ci.yml`](.github/workflows/ci.yml)) menjalankan:

- `go vet`
- `staticcheck`
- `go test ./... -race -coverprofile=coverage.out`
- **Threshold minimum: 70%** total statement coverage

Jalankan perintah coverage di atas secara lokal sebelum push untuk memastikan CI lulus.

---

## CI/CD

Pipeline GitHub Actions pada setiap push/PR ke `main`:

1. **Lint & Test** — `go vet`, `staticcheck`, unit test + coverage
2. **Docker** (hanya push ke `main`) — build & push image ke `ghcr.io/<owner>/hexa-go`

---

## Referensi

- [Wiki Dokumentasi](https://github.com/rulzi/hexa-go/wiki)
- [Hexagonal Architecture — Alistair Cockburn](https://alistair.cockburn.us/hexagonal-architecture/)
- [Clean Architecture — Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

## License

MIT License

## Author

**Khoirul Afandi**

- Instagram: [@afandi_](https://instagram.com/afandi_)
- LinkedIn: [Khoirul Afandi](https://www.linkedin.com/in/khoirulafandi/)
