# 🌀 Unbound Base v1.0 — Barebone Social Backend

**Unbound Base** adalah *skeleton backend* dari proyek sosial media terdistribusi **Unbound**, dibangun dengan **Go + Fiber + PostgreSQL**.  
Versi ini merupakan rilis **v1.0 stable**, mencakup semua fondasi utama untuk sistem sosial: autentikasi, posting, komentar, like, follow, feed, search, dan notifikasi.

---

## ✨ What's New in v1.0
- ✅ **Full Auth System** — Register, login, refresh token, dan logout  
- ✅ **Notification System** — Fetch & mark as read untuk like, comment, dan follow  
- ✅ **Feed System** — Timeline publik & following dengan pagination dan sorting  
- ✅ **User & Follow System** — Profil, follow/unfollow, list followers/following  
- ✅ **Post & Comment System** — CRUD post, komentar, likes, dan counting  
- ✅ **Search System** — Pencarian user & post dengan filter dan sort  

---
## 🚧 Ongoing Development
- 💬 Direct Message / Chat System — Sistem chat antar user (private messaging)
- 🌀 Topic - Sistem pengelompokan makna postingan
- 🌐 Realtime Update — WebSocket layer untuk notifikasi & chat
- 🧩 Microservice Split — Pisahkan auth, post, user, dan notification ke service mandiri
- 🐳 Docker Compose Setup — Containerisasi full stack backend
- 🔍 ElasticSearch Integration — Pencarian lebih cepat dan relevan
- ⚙️ Machine Learning Integration - Rekomendasi feed yang lebih relevan

---
## ⚙️ Tech Stack
- **Go (Fiber v2)** – Fast HTTP framework  
- **GORM + PostgreSQL** – ORM dan database utama  
- **JWT (golang-jwt/v5)** – Autentikasi stateless (access + refresh token)  
- **Notification System** – Event-based alert untuk like, comment, follow  
- **Docker (future)** – Containerization  
- **Kafka, MinIO, Redis (planned)** – Event bus, storage, caching  

---

## 📡 API Endpoints

| Method | Endpoint | Deskripsi |
|:--|:--|:--|
| `POST` | `/auth/register` | Register user baru |
| `POST` | `/auth/login` | Login dan dapatkan JWT |
| `POST` | `/auth/refresh` | Refresh access token |
| `POST` | `/auth/logout` | Logout dan hapus refresh token |
| `POST` | `/posts` | Buat posting (auth) |
| `PUT` | `/posts/:id` | Edit posting milik sendiri |
| `DELETE` | `/posts/:id` | Hapus posting milik sendiri |
| `GET` | `/feed` | Lihat timeline publik |
| `GET` | `/feed/following` | Lihat timeline dari user yang di-follow |
| `GET` | `/users/:username` | Lihat profil dan post user |
| `POST` | `/users/:username/follow` | Follow / Unfollow user |
| `GET` | `/users/:username/followers` | Lihat siapa saja followers user |
| `GET` | `/users/:username/following` | Lihat siapa saja yang di-follow user |
| `POST` | `/posts/:id/like` | Like / Unlike post |
| `GET` | `/posts/:id/likes` | Hitung total likes |
| `POST` | `/posts/:id/comments` | Tambah komentar ke post |
| `PUT` | `/posts/:post_id/comments/:id` | Edit komentar milik sendiri |
| `GET` | `/posts/:id/comments` | Lihat semua komentar |
| `DELETE` | `/posts/:post_id/comments/:id` | Hapus komentar milik sendiri |
| `POST` | `/search?query=` | Pencarian beserta filter by user,post,oldest/newest |
| `GET` | `/notifications` | Ambil daftar notifikasi (like, comment, follow) |
| `POST` | `/notifications/read` | Tandai semua notifikasi user sebagai dibaca |

---

## 🧱 Struktur
```
unbound/
├── cmd/server/           # Entry point
├── internal/
│   ├── auth/             # Register, login, JWT, refresh, logout
│   ├── post/             # Post, like, comment, feed, edit
│   ├── user/             # Profile & follow system
│   ├── search/           # Pencarian user & post
│   ├── notification/     # Sistem notifikasi
│   └── common/           # DB, middleware, utils
└── go.mod
```

---

## 🚀 Jalankan Lokal
```bash
# clone repo
git clone https://github.com/azdharsyahputra/unbound-base.git
cd unbound-base

# buat .env
DB_HOST=localhost
DB_USER=postgres
DB_PASS=<password>
DB_NAME=unbound_db
DB_PORT=5432
JWT_SECRET=dev-secret

# jalanin server
go run cmd/server/main.go
```
Server jalan di: **http://localhost:8080**

---

## 🧑‍💻 Author
**Muhammad Azdhar Syahputra (Ajar)**  
Backend Systems & Architecture — [@azdharsyahputra](https://github.com/azdharsyahputra)