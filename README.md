# 🌀 Unbound Base v1.1 — Private Chat & Notifications Realtime

**Unbound Base** adalah *barebone backend* dari proyek sosial media terdistribusi **Unbound**, dibangun dengan **Go + Fiber + PostgreSQL**.  
Versi **v1.1** memperkenalkan **Direct Messaging (Private Chat)** dengan dukungan **realtime WebSocket** dan **notifikasi otomatis**.

---

## ✨ What's New in v1.1
- 💬 **Private Chat System** — Chat pribadi antar user dengan tabel `chats` & `messages`
- ⚡ **Realtime WebSocket Layer** — Komunikasi dua arah instan via endpoint `/ws/chat/:chat_id`
- 🔔 **Auto Notification Sync** — Pesan baru via WebSocket otomatis membuat notifikasi untuk penerima
- 👁️ **Read & Delivery Receipt** — Pesan ditandai *delivered* ketika lawan bicara terhubung dan bisa ditandai *read* via API
- 🧩 **Full Auth Integration** — Semua endpoint chat dilindungi JWT middleware
- 🪶 **Backward Compatible** — Semua API dari v1.0 masih berfungsi penuh

---

## ⚙️ Tech Stack
- **Go (Fiber v2)** – Fast HTTP framework  
- **GORM + PostgreSQL** – ORM dan database utama  
- **JWT (golang-jwt/v5)** – Autentikasi stateless  
- **WebSocket (fiber/contrib)** – Realtime communication layer  
- **Notification System** – Event-based alert untuk like, comment, follow, dan chat  
- **Docker (planned)** – Containerization  
- **Redis / Kafka (future)** – Event bus & caching layer  

---

## 📡 API Endpoints

### 🔐 Auth
| Method | Endpoint | Deskripsi |
|:--|:--|:--|
| `POST` | `/auth/register` | Register user baru |
| `POST` | `/auth/login` | Login dan dapatkan JWT |
| `POST` | `/auth/refresh` | Refresh access token |
| `POST` | `/auth/logout` | Logout dan hapus refresh token |

### 📰 Post & Feed
| Method | Endpoint | Deskripsi |
|:--|:--|:--|
| `POST` | `/posts` | Buat posting (auth) |
| `PUT` | `/posts/:id` | Edit posting milik sendiri |
| `DELETE` | `/posts/:id` | Hapus posting milik sendiri |
| `GET` | `/feed` | Timeline publik |
| `GET` | `/feed/following` | Timeline dari user yang di-follow |

### 👥 User & Follow
| Method | Endpoint | Deskripsi |
|:--|:--|:--|
| `GET` | `/users/:username` | Lihat profil user |
| `POST` | `/users/:username/follow` | Follow / Unfollow user |
| `GET` | `/users/:username/followers` | Lihat followers |
| `GET` | `/users/:username/following` | Lihat yang di-follow |

### 💬 Chat & Messages
| Method | Endpoint | Deskripsi |
|:--|:--|:--|
| `GET` | `/chats` | Ambil daftar chat user login |
| `POST` | `/chats/:user_id` | Buat atau ambil chat dengan user tertentu |
| `GET` | `/chats/:chat_id/messages` | Ambil semua pesan dalam chat |
| `POST` | `/chats/:chat_id/messages` | Kirim pesan baru |
| `PUT` | `/chats/:chat_id/read` | Tandai semua pesan sebagai dibaca |
| `GET` | `/ws/chat/:chat_id?token=` | Realtime WebSocket endpoint |

### 🔔 Notifications
| Method | Endpoint | Deskripsi |
|:--|:--|:--|
| `GET` | `/notifications` | Ambil semua notifikasi user |
| `POST` | `/notifications/read` | Tandai semua notifikasi sebagai dibaca |

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
│   ├── chat/             # Private chat, WebSocket, message delivery
│   ├── notification/     # Sistem notifikasi (event-based)
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