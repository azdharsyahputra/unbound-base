# 🌀 Unbound Base — Barebone Backend

**Unbound Base** adalah *skeleton backend* dari proyek sosial media terdistribusi **Unbound**, dibangun dengan **Go + Fiber + PostgreSQL**.  
Tujuan repo ini adalah menyediakan fondasi API utama sebelum dipisah menjadi microservices.

---

## ⚙️ Tech Stack
- **Go (Fiber v2)** – Fast HTTP framework  
- **GORM + PostgreSQL** – ORM dan database utama  
- **JWT (golang-jwt/v5)** – Autentikasi stateless  
- **Docker (future)** – Containerization  
- **Kafka, MinIO, Redis (planned)** – Event bus, storage, caching  

---

## 📡 API Endpoints

| Method | Endpoint | Deskripsi |
|:--|:--|:--|
| `POST` | `/auth/register` | Register user baru |
| `POST` | `/auth/login` | Login dan dapatkan JWT |
| `POST` | `/posts` | Buat posting (auth) |
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
| `GET` | `/posts/:id/comments` | Lihat semua komentar |
| `DELETE` | `/posts/:post_id/comments/:id` | Hapus komentar milik sendiri |
---

## 🧱 Struktur
```
unbound/
├── cmd/server/           # Entry point
├── internal/
│   ├── auth/             # Register, login, JWT
│   ├── post/             # Post, like, comment, feed
│   ├── user/             # Profile endpoint
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