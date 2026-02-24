# 🎉 Event Invitation — Go Backend

Free digital invitation platform for weddings, birthdays, communities, and other events.

## Tech Stack

- **Language**: Go 1.22
- **Framework**: Gin
- **Database**: PostgreSQL (sqlx)
- **Cache**: Redis
- **Auth**: JWT (golang-jwt/jwt)
- **Password**: bcrypt

---

## Arsitektur

```
Request → Handler → Service → Repository → Database
```

Layer yang jelas dan terpisah:
- **Handler**: Parse request, call service, return response
- **Service**: Business logic (validasi, slug generation, ownership check)
- **Repository**: Query database saja
- **Domain**: Struct + interface kontrak


## Quick Start

### 1. Siapkan environment
```bash
cp .env.example .env
# Edit .env sesuai kebutuhan
```

### 2. Jalankan dengan Docker
```bash
docker-compose up -d
```
Database akan otomatis ter-migrate saat PostgreSQL pertama kali start.

### 3. Atau jalankan manual
```bash
# Pastikan PostgreSQL & Redis sudah jalan
psql -U postgres -d event_invitation -f migrations/0001_init_schema.up.sql

go mod tidy
go run ./cmd/api
```

Server berjalan di `http://localhost:8080`

---

## API Endpoints

### Auth (Public)
| Method | Endpoint | Keterangan |
|--------|----------|------------|
| POST | `/api/v1/auth/register` | Daftar akun |
| POST | `/api/v1/auth/login` | Login, dapat JWT token |

### Templates (Public)
| Method | Endpoint | Keterangan |
|--------|----------|------------|
| GET | `/api/v1/templates` | List semua template (filter: `?category=wedding`) |
| GET | `/api/v1/templates/:id` | Detail template + sections |

### Public Event
| Method | Endpoint | Keterangan |
|--------|----------|------------|
| GET | `/api/v1/e/:slug` | Halaman undangan publik |
| POST | `/api/v1/events/:id/rsvp` | Submit RSVP (publik) |

### Events (🔒 JWT Required)
| Method | Endpoint | Keterangan |
|--------|----------|------------|
| POST | `/api/v1/events` | Buat event baru |
| GET | `/api/v1/events` | List event milik user |
| GET | `/api/v1/events/:id` | Detail event |
| PATCH | `/api/v1/events/:id` | Update event |
| DELETE | `/api/v1/events/:id` | Hapus event |
| PATCH | `/api/v1/events/:id/publish` | Publish/unpublish |
| PUT | `/api/v1/events/:id/theme` | Update tema (warna, font, dll) |
| PATCH | `/api/v1/events/:id/sections/:sectionId` | Update konten section |
| GET | `/api/v1/events/:id/guests` | Daftar tamu RSVP |
| POST | `/api/v1/events/:id/media` | Upload gambar/video/audio |
| GET | `/api/v1/events/:id/media` | List media event |
| DELETE | `/api/v1/events/:id/media/:mediaId` | Hapus media |

---

## Environment Variables

| Key | Default | Keterangan |
|-----|---------|------------|
| `APP_PORT` | `8080` | Port server |
| `DB_HOST` | `localhost` | PostgreSQL host |
| `DB_NAME` | `event_invitation` | Nama database |
| `JWT_SECRET` | — | Secret untuk JWT (ganti di production!) |
| `JWT_EXPIRY_HOURS` | `72` | Masa berlaku token (jam) |
| `STORAGE_BASE_PATH` | `./uploads` | Folder penyimpanan file upload |
| `STORAGE_BASE_URL` | `http://localhost:8080/uploads` | Base URL untuk akses file |
