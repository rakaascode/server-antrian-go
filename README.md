# 🚀 Server Antrian Bengkel Lautan Teduh (Go Backend API)

[![Go](https://img.shields.io/badge/Language-Go-blue.svg)](https://go.dev)
[![Gin](https://img.shields.io/badge/Framework-Gin-green.svg)](https://github.com/gin-gonic/gin)
[![GORM](https://img.shields.io/badge/ORM-GORM-orange.svg)](https://gorm.io)
[![PostgreSQL](https://img.shields.io/badge/Database-PostgreSQL-blue.svg)](https://www.postgresql.org)
[![Docker](https://img.shields.io/badge/Container-Docker-2496ED.svg)](https://www.docker.com)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

> Backend API production-ready untuk sistem antrian bengkel dengan autentikasi Google & Admin, dibangun menggunakan Golang dengan arsitektur clean.

---

## 🎯 Overview

Project ini merupakan backend service untuk sistem antrian bengkel yang menyediakan:

* 🔐 **Authentication System**

  * Login User via Google Token
  * Login Admin via Username & Password
* 📦 **REST API CRUD**
* 🧱 **Clean Architecture (Handler → Service → Repository)**
* 🐳 **Dockerized Environment**

---

## ✨ Features

### 🔐 Authentication

* 🔵 Google Login (OAuth token based)
* 🔴 Admin Login (username & password)
* 🔑 JWT Token generation
* 🔒 Secure token verification

---

### 📊 User & Admin Management

* CRUD User
* CRUD Admin
* Role separation (User vs Admin)

---

### ⚙️ Technical Features

* ⚡ High performance with Golang
* 🗄️ PostgreSQL database
* 🧩 GORM ORM
* 🌐 RESTful API (Gin)
* 🐳 Docker Compose setup
* 🧱 Layered architecture

---

## 🏗️ Architecture

### 🔄 Request Flow

```id="archflow"
Request
  ↓
Handler (HTTP Layer)
  ↓
Service (Business Logic)
  ↓
Repository (Database Access)
  ↓
PostgreSQL
```

---

## 📁 Struktur Project

```
server-antrian-go/
│
├── cmd/
│   └── main.go              # Entry point aplikasi
│
├── config/
│   └── database.go          # Koneksi database (GORM)
│
├── models/
│   ├── user.go              # Model user
│   └── admin.go             # Model admin
│
├── repository/
│   ├── user_repository.go   # Query user ke database
│   └── admin_repository.go  # Query admin ke database
│
├── service/
│   ├── auth_service.go      # 🔥 Logic autentikasi (Google & Admin)
│   ├── user_service.go
│   └── admin_service.go
│
├── handler/
│   ├── auth_handler.go      # 🔥 Endpoint login
│   ├── user_handler.go
│   └── admin_handler.go
│
├── routes/
│   └── routes.go            # Routing API
│
├── utils/
│   ├── jwt.go               # Generate JWT
│   └── google.go            # Verifikasi token Google
│
└── .env                     # Konfigurasi environment
```
---

## 🐳 Docker Setup

### 📦 Services

* **app** → Backend Go
* **db** → PostgreSQL 15

---

### ▶️ Run Application

```bash
docker-compose up --build
```

---

### 🌐 Access API

```id="url"
http://localhost:8080
```

---

## ⚙️ Environment Configuration

Sudah di-handle oleh Docker Compose

---

## 🔐 Authentication Flow

### 🔵 Google Login

```id="googleflow"
Client → kirim Google Token
        ↓
Backend verify ke Google
        ↓
User dibuat / ditemukan
        ↓
JWT dikirim ke client
```

---

### 🔴 Admin Login

```id="adminflow"
Client → username & password
        ↓
Cek ke database
        ↓
Generate JWT
        ↓
Return token
```

---

## 📌 API Endpoints

### 🔐 Auth

```id="authapi"
POST /api/login/google
POST /api/login/admin
```

---

### 👤 User

```id="userapi"
GET    /api/users
POST   /api/users
GET    /api/users/:id
PUT    /api/users/:id
DELETE /api/users/:id
```

---

### 👨‍💼 Admin

```id="adminapi"
GET    /api/admins
POST   /api/admins
```

---

## 🔧 Key Implementation

### 🔐 JWT Generation

* Token dibuat di `utils/jwt.go`
* Digunakan untuk autentikasi request selanjutnya

---

### 🌐 Google Token Verification

* Token diverifikasi di `utils/google.go`
* Tidak percaya token langsung dari client

---

## ⚠️ Best Practices

* 🔒 Jangan commit `.env`
* 🔑 Gunakan bcrypt untuk password admin
* 🔐 Selalu verify Google token
* 🧪 Tambahkan validation di handler

---

## 🚀 Roadmap

* [ ] JWT Middleware
* [ ] Role-based access control
* [ ] Pagination
* [ ] Logging system
* [ ] Unit testing
* [ ] API documentation (Swagger)

---

## 🐛 Troubleshooting

### ❌ Cannot connect to DB

* Pastikan container `db` sudah running
* Cek `DB_HOST=db` (bukan localhost)

---

### ❌ Port already in use

* Ganti port di docker-compose

---

### ❌ Google login gagal

* Pastikan `GOOGLE_CLIENT_ID` benar
* Token belum expired

---

## 📄 License

Licensed under the **Apache License 2.0**.

---

## 🙌 Acknowledgments

* Golang community
* GORM contributors
* Gin framework
* PostgreSQL ecosystem
* Docker community
