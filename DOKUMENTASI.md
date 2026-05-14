# 📋 API Dokumentasi — Lautan Teduh Antrian System

---

## 1. 🌐 Overview

API untuk sistem manajemen antrian bengkel motor **Lautan Teduh** yang mendukung multi-cabang.

**Fitur utama:**
- Manajemen antrian per cabang secara real-time
- Autentikasi multi-peran: User Android (Google), Admin Cabang
- **Profil user otomatis** dari Google (nama, email, foto avatar) — tersimpan di database
- **Riwayat antrian user** beserta info cabang (nama, alamat, kota) yang bisa diakses via JWT
- Notifikasi pengingat otomatis via WhatsApp (Fonnte)
- In-app broadcast: promo & pengumuman per cabang di aplikasi Android
- CRM pengiriman pesan WA manual atau dari data antrian

> 📌 **Dua Channel Notifikasi yang Berbeda:**
> - **WhatsApp** (`/crm/*`): Dikirim ke nomor WA pelanggan via Fonnte — langsung masuk WA
> - **In-App Broadcast** (`/broadcast`): Disimpan di database — user baca di inbox aplikasi Android

| Info | Detail |
|---|---|
| **Base URL** | `http://localhost:8080/api` |
| **Format** | JSON |
| **Auth** | JWT Bearer Token |

---

## 2. 🔑 Authentication

API ini menggunakan JWT (JSON Web Token). Setelah login, sertakan token di setiap request yang membutuhkan autentikasi.

**Header:**
```
Authorization: Bearer <token>
```

**Jenis Akun:**

| Role | Login Via | Akses |
|---|---|---|
| `user` | **Google Sign-In** (Android) | Ambil antrian, lihat antrian sendiri |
| `admin` | **Username + Password** | Kelola semua antrian di cabangnya |

---

## 3. 🔗 Daftar Endpoint

| Method | Endpoint | Auth | Deskripsi |
|---|---|---|---|
| POST | `/auth/google` | ❌ | Login user via Google Sign-In (Android) |
| POST | `/auth/admin/login` | ❌ | Login admin cabang (username + password) |
| GET | `/cabang` | ❌ | List semua cabang |
| GET | `/cabang/:id` | ❌ | Detail satu cabang |
| POST | `/cabang` | ✅ Admin | Buat cabang baru |
| PUT | `/cabang/:id` | ✅ Admin | Update data cabang |
| DELETE | `/cabang/:id` | ✅ Admin | Hapus cabang |
| GET | `/cabang/:id/antrian` | ❌ | List antrian cabang (publik, tanpa data sensitif) |
| GET | `/cabang/:id/antrian/status` | ❌ | 🔴 Nomor yang sedang dipanggil + total menunggu |
| GET | `/cabang/:id/antrian/detail` | ✅ Admin | List antrian cabang (detail lengkap) |
| POST | `/antrian` | ✅ User | Ambil nomor antrian |
| GET | `/antrian/me` | ✅ User | Semua antrian milik saya |
| GET | `/antrian/:id` | ✅ User/Admin | Detail antrian (owner / admin cabang) |
| GET | `/antrian/:id/posisi` | ✅ User | Posisi di antrian & nomor yang sedang dilayani |
| POST | `/antrian/call-next` | ✅ Admin | Panggil antrian berikutnya |
| PUT | `/antrian/:id/selesai` | ✅ Admin | Tandai antrian selesai |
| DELETE | `/antrian/:id` | ✅ Admin | Hapus/batalkan antrian |
| **GET** | **`/user/profile`** | ✅ User | **Profil user (nama, email, avatar) + riwayat antrian & cabang** |
| **PUT** | **`/user/profile`** | ✅ User | **Update profil user (alamat, kota, provinsi, kode pos, promo_aktif, avatar)** |
| **GET** | **`/users/kontak`** | ✅ User | **Lihat nomor WA yang tersimpan** |
| **POST** | **`/users/kontak`** | ✅ User | **Simpan nomor WA user** |
| **PUT** | **`/users/kontak`** | ✅ User | **Update nomor WA user** |
| **DELETE** | **`/users/kontak`** | ✅ User | **Hapus nomor WA user** |
| POST | `/crm/send` | ✅ Admin | Kirim WA manual ke nomor tertentu |
| POST | `/crm/reminders` | ✅ Admin | Kirim pengingat WA dari data antrian |
| POST | `/broadcast` | ✅ Admin | Kirim broadcast in-app (promo / per cabang) |
| GET | `/broadcast/all` | ✅ Admin | Lihat semua broadcast yang pernah dikirim |
| GET | `/broadcast` | ✅ User | Inbox notifikasi (promo + antrian cabang saya) |
| GET | `/broadcast/:id` | ✅ User | Detail lengkap satu notifikasi |
| **GET** | **`/admin/users/kontak`** | ✅ Admin | **Ambil semua nama dan nomor WA user** |
| GET | `/users` | ✅ Admin | List semua user |
| GET | `/users/:id` | ✅ Admin | Detail satu user |
| POST | `/users` | ✅ Admin | Buat user/admin baru |
| PUT | `/users/:id` | ✅ Admin | Update user |
| DELETE | `/users/:id` | ✅ Admin | Hapus user |
| POST | `/super/admin/cabang` | ✅ Super Admin | Buat admin cabang baru |
| GET | `/super/admins` | ✅ Super Admin | Lihat seluruh admin |
| GET | `/super/cabang/:id/admins` | ✅ Super Admin | Lihat admin di cabang tertentu |
| PUT | `/super/admins/:id/assign` | ✅ Super Admin | Assign cabang ke admin |
| DELETE | `/super/admins/:id/assign` | ✅ Super Admin | Melepas admin dari cabangnya |

---

## 4. 📖 Detail Endpoint

---

### AUTH

---


#### `POST /auth/google`
Login user Android menggunakan Google ID Token.

**📥 Request Body:**
```json
{
  "id_token": "eyJhbGciOiJSUzI1NiIs..."
}
```

| Field | Tipe | Required | Keterangan |
|---|---|---|---|
| `id_token` | string | ✅ | ID Token dari Google Sign-In SDK Android |

> Jika akun belum ada, sistem **otomatis membuat akun baru** dari data Google (nama, email, foto profil).
> Jika akun sudah ada dan foto profil Google berubah, `avatar_url` di database **otomatis diperbarui** setiap login.

**📤 Response `200`:**
```json
{
  "success": true,
  "message": "Login Google berhasil",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": 2,
      "name": "Budi Google",
      "email": "budi@gmail.com",
      "avatar_url": "https://lh3.googleusercontent.com/a/ACg8ocIx...",
      "role": "user"
    }
  }
}
```

> **`avatar_url`** adalah link foto profil Google (CDN Google). Langsung bisa ditampilkan di `<img src="...">` pada UI Android.

---

#### `POST /auth/admin/login`
Login admin cabang dengan username dan password.

**📥 Request Body:**
```json
{
  "username": "admin_kedaton",
  "password": "password123"
}
```

| Field | Tipe | Required |
|---|---|---|
| `username` | string | ✅ |
| `password` | string | ✅ |

**📤 Response `200`:**
```json
{
  "success": true,
  "message": "Login admin berhasil",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": 10,
      "name": "Admin Kedaton",
      "username": "admin_kedaton",
      "role": "admin",
      "cabang_id": 1
    }
  }
}
```

**🧪 cURL:**
```bash
curl -X POST http://localhost:8080/api/auth/admin/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin_kedaton","password":"password123"}'
```

---

### CABANG

---

#### `GET /cabang`
Ambil daftar semua cabang. Tidak memerlukan login.

**📤 Response `200`:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "nama": "Lautan Teduh Kedaton",
      "alamat": "Jl. Teuku Umar No.15D, Kedaton",
      "kota": "Bandar Lampung",
      "no_telp": "081367846069",
      "latitude": -5.3795,
      "longitude": 105.261
    }
  ]
}
```

**🧪 cURL:**
```bash
curl http://localhost:8080/api/cabang
```

---

#### `GET /cabang/:id`
Ambil detail satu cabang.

**📤 Response `200`:** *(struktur sama seperti di atas, satu objek)*

**📤 Response `404`:**
```json
{ "success": false, "message": "Cabang tidak ditemukan" }
```

---

#### `POST /cabang` — Admin
Buat cabang baru.

**📥 Request Body:**
```json
{
  "nama": "Lautan Teduh Kedaton",
  "alamat": "Jl. Teuku Umar No.15D, Kedaton",
  "kota": "Bandar Lampung",
  "no_telp": "081367846069",
  "latitude": -5.3795,
  "longitude": 105.2610
}
```

| Field | Tipe | Required |
|---|---|---|
| `nama` | string | ✅ |
| `alamat` | string | ✅ |
| `kota` | string | ✅ |
| `no_telp` | string | ❌ |
| `latitude` | float | ❌ |
| `longitude` | float | ❌ |

**📤 Response `201`:**
```json
{ "success": true, "message": "Cabang berhasil dibuat", "data": { ... } }
```

---

### ANTRIAN

---

#### `POST /antrian` — User (wajib login)
Ambil nomor antrian di cabang tertentu. Antrian otomatis terhubung ke akun user.

**📥 Request Body:**
```json
{
  "cabang_id": 1,
  "nama_pemilik": "Budi Santoso",
  "no_hp": "081234567890",
  "merk_motor": "Honda",
  "tipe_motor": "Vario 150",
  "no_rangka": "MH1JM2109NK000123",
  "no_mesin": "JM21E1000123",
  "tahun_pembuatan": 2021,
  "tanggal_kedatangan": "2024-11-20T00:00:00Z",
  "estimasi_jam": "09:00",
  "catatan": "Servis rutin dan ganti oli",
  "reminder_aktif": true,
  "no_wa_reminder": "081234567890"
}
```

| Field | Tipe | Required | Keterangan |
|---|---|---|---|
| `cabang_id` | int | ✅ | ID cabang tujuan |
| `nama_pemilik` | string | ✅ | Nama sesuai STNK |
| `no_hp` | string | ✅ | No HP sesuai STNK |
| `merk_motor` | string | ✅ | Misal: Honda, Yamaha |
| `tipe_motor` | string | ✅ | Misal: Vario 150, NMAX |
| `no_rangka` | string | ✅ | Nomor rangka kendaraan |
| `no_mesin` | string | ✅ | Nomor mesin kendaraan |
| `tahun_pembuatan` | int | ✅ | Misal: 2021 |
| `tanggal_kedatangan` | datetime | ✅ | Format ISO 8601 |
| `estimasi_jam` | string | ✅ | Format: `"09:00"` |
| `catatan` | string | ❌ | Keluhan / catatan tambahan |
| `reminder_aktif` | bool | ❌ | `true` = aktifkan pengingat WA |
| `no_wa_reminder` | string | ⚠️ | Wajib jika `reminder_aktif: true` |

**📤 Response `201`:**
```json
{
  "success": true,
  "message": "Antrian berhasil diambil",
  "data": {
    "id": 5,
    "cabang_id": 1,
    "user_id": 2,
    "nomor_antrian": 5,
    "status": "menunggu",
    "nama_pemilik": "Budi Santoso",
    "no_hp": "081234567890",
    "merk_motor": "Honda",
    "tipe_motor": "Vario 150",
    "no_rangka": "MH1JM2109NK000123",
    "no_mesin": "JM21E1000123",
    "tahun_pembuatan": 2021,
    "tanggal_kedatangan": "2024-11-20T00:00:00Z",
    "estimasi_jam": "09:00",
    "reminder_aktif": true,
    "no_wa_reminder": "081234567890",
    "created_at": "2024-11-19T08:00:00Z"
  }
}
```

**🧪 cURL:**
```bash
curl -X POST http://localhost:8080/api/antrian \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token_user>" \
  -d '{
    "cabang_id": 1,
    "nama_pemilik": "Budi Santoso",
    "no_hp": "081234567890",
    "merk_motor": "Honda",
    "tipe_motor": "Vario 150",
    "no_rangka": "MH1JM2109NK000123",
    "no_mesin": "JM21E1000123",
    "tahun_pembuatan": 2021,
    "tanggal_kedatangan": "2024-11-20T00:00:00Z",
    "estimasi_jam": "09:00",
    "reminder_aktif": true,
    "no_wa_reminder": "081234567890"
  }'
```

---

#### `GET /cabang/:id/antrian`
Lihat antrian di cabang tertentu. **Tidak memerlukan login.** Data sensitif (No HP, No Rangka, dll) **tidak ditampilkan**.

**Query Params:** `?status=menunggu` *(opsional, filter: menunggu / dipanggil / selesai)*

**🧪 cURL:**
```bash
curl "http://localhost:8080/api/cabang/1/antrian?status=menunggu"
```

**📤 Response `200`:**
```json
{
  "success": true,
  "total_antrian": 3,
  "data": [
    {
      "id": 1,
      "nomor_antrian": 1,
      "status": "menunggu",
      "estimasi_jam": "09:00",
      "tanggal_kedatangan": "2024-11-20T00:00:00Z"
    }
  ]
}
```

---

#### `GET /cabang/:id/antrian/status` 🔴 — Public
Info **realtime** antrian di cabang: nomor yang **sedang dipanggil** sekarang + total antrian yang masih menunggu. Cocok untuk ditampilkan di layar display bengkel atau aplikasi Android.

**🧪 cURL:**
```bash
curl http://localhost:8080/api/cabang/1/antrian/status
```

**📤 Response `200` (ada yang dipanggil):**
```json
{
  "success": true,
  "data": {
    "nomor_dipanggil": 5,
    "status_panggil": "dipanggil",
    "total_menunggu": 8
  }
}
```

**📤 Response `200` (belum ada yang dipanggil):**
```json
{
  "success": true,
  "data": {
    "nomor_dipanggil": null,
    "status_panggil": "belum ada",
    "total_menunggu": 12
  }
}
```

---

#### `GET /antrian/me` — User
Lihat semua antrian yang pernah dibuat oleh user yang sedang login.

**🧪 cURL:**
```bash
curl http://localhost:8080/api/antrian/me \
  -H "Authorization: Bearer <token_user>"
```

---

#### `GET /user/profile` — User (wajib login)
Ambil profil lengkap user yang sedang login berdasarkan JWT. Response mencakup:
- Data profil: nama, email, foto avatar Google
- Riwayat seluruh antrian yang pernah diambil user, **lengkap dengan info cabang** (nama, alamat, kota, no telp)

Tidak perlu kirim ID user — server otomatis baca dari JWT token.

**🧪 cURL:**
```bash
curl http://localhost:8080/api/user/profile \
  -H "Authorization: Bearer <token_user>"
```

**📤 Response `200`:**
```json
{
  "success": true,
  "data": {
    "id": 2,
    "name": "Budi Santoso",
    "email": "budi@gmail.com",
    "avatar_url": "https://lh3.googleusercontent.com/a/ACg8ocIx...",
    "role": "user",
    "created_at": "2024-11-01T08:00:00Z",
    "antrian": [
      {
        "id": 5,
        "nomor_antrian": 3,
        "status": "selesai",
        "tanggal_kedatangan": "2024-11-20T00:00:00Z",
        "estimasi_jam": "09:00",
        "merk_motor": "Honda",
        "tipe_motor": "Vario 150",
        "created_at": "2024-11-19T08:00:00Z",
        "cabang": {
          "id": 1,
          "nama": "Lautan Teduh Kedaton",
          "alamat": "Jl. Teuku Umar No.15D, Kedaton",
          "kota": "Bandar Lampung",
          "no_telp": "081367846069"
        }
      },
      {
        "id": 3,
        "nomor_antrian": 7,
        "status": "selesai",
        "tanggal_kedatangan": "2024-10-05T00:00:00Z",
        "estimasi_jam": "10:30",
        "merk_motor": "Yamaha",
        "tipe_motor": "NMAX",
        "created_at": "2024-10-04T09:00:00Z",
        "cabang": {
          "id": 3,
          "nama": "Lautan Teduh Tirtayasa",
          "alamat": "Jl. Tirtayasa No.12",
          "kota": "Bandar Lampung",
          "no_telp": "081234000000"
        }
      }
    ]
  }
}
```

| Field | Keterangan |
|---|---|
| `avatar_url` | Link foto profil Google. Langsung dipakai sebagai `src` di komponen gambar Android |
| `antrian[]` | Riwayat antrian diurutkan dari **terbaru ke terlama** |
| `antrian[].cabang` | Info cabang tempat antrian diambil (nama, alamat, kota, no_telp) |
| `antrian[].status` | `menunggu` / `dipanggil` / `selesai` |

**📤 Response `404`:**
```json
{ "success": false, "message": "User tidak ditemukan" }
```

---

#### `PUT /user/profile` — User (wajib login)
Memperbarui informasi profil pengguna seperti nama, alamat, dan persetujuan menerima info promo dari WhatsApp. Alias endpoint ini: `PUT /users/profile`.

**📥 Request Body:**
```json
{
  "name": "Budi Santoso",
  "avatar_url": "https://url-bucket-avatar.com/budi.jpg",
  "alamat": "Jl. Mawar Merah No. 12",
  "kota": "Surabaya",
  "provinsi": "Jawa Timur",
  "kode_pos": "60241",
  "promo_aktif": true
}
```

Semua parameter di atas **opsional**. Kirim hanya *field* yang ingin diperbarui, misalnya hanya `"promo_aktif": true`.

**📤 Response `200`:**
```json
{
  "success": true,
  "message": "Profil berhasil diperbarui",
  "data": { ... }
}
```

---

#### `GET /antrian/:id/posisi` — User
Cek posisi antrian user: **berapa orang di depan** dan **nomor yang sedang dilayani**. Berguna untuk notifikasi giliran di aplikasi Android.

Header: `Authorization: Bearer <token_user>`

> Hanya bisa diakses oleh owner antrian tersebut atau admin cabang yang sama.

**🧪 cURL:**
```bash
curl http://localhost:8080/api/antrian/12/posisi \
  -H "Authorization: Bearer <token_user>"
```

**📤 Response `200` (masih menunggu):**
```json
{
  "success": true,
  "data": {
    "nomor_antrian": 8,
    "status": "menunggu",
    "posisi": 2,
    "nomor_dipanggil": 5,
    "pesan": "2 orang di depan Anda."
  }
}
```

**📤 Response `200` (giliran berikutnya):**
```json
{
  "success": true,
  "data": {
    "nomor_antrian": 6,
    "status": "menunggu",
    "posisi": 0,
    "nomor_dipanggil": 5,
    "pesan": "Anda adalah antrian berikutnya!"
  }
}
```

**📤 Response `200` (sudah dipanggil):**
```json
{
  "success": true,
  "data": {
    "nomor_antrian": 6,
    "status": "dipanggil",
    "posisi": 0,
    "nomor_dipanggil": 6,
    "pesan": "Nomor antrian Anda sedang dipanggil! Segera ke loket."
  }
}
```

| `posisi` | Artinya |
|---|---|
| `0` + status `menunggu` | Giliran berikutnya |
| `0` + status `dipanggil` | Sedang dipanggil sekarang |
| `> 0` | Masih ada `posisi` orang di depan |

---


#### `GET /antrian/:id` — User / Admin
Lihat detail lengkap satu antrian. Hanya bisa diakses oleh:
- User yang membuat antrian tersebut
- Admin dari cabang yang sama

**📤 Response `403` (jika bukan owner/admin):**
```json
{ "success": false, "message": "Akses ditolak" }
```

---

#### `POST /antrian/call-next` — Admin
Panggil antrian berikutnya (status: menunggu → dipanggil). Jika antrian memiliki `reminder_aktif: true`, sistem **otomatis mengirim WA**.

**🧪 cURL:**
```bash
curl -X POST http://localhost:8080/api/antrian/call-next \
  -H "Authorization: Bearer <token_admin>"
```

**📤 Response `200`:**
```json
{
  "success": true,
  "message": "Antrian dipanggil",
  "data": { "id": 1, "nomor_antrian": 1, "status": "dipanggil", ... }
}
```

---

#### `PUT /antrian/:id/selesai` — Admin
Tandai antrian selesai (status: dipanggil → selesai). Jika `reminder_aktif: true`, sistem kirim WA notifikasi selesai.

**🧪 cURL:**
```bash
curl -X PUT http://localhost:8080/api/antrian/1/selesai \
  -H "Authorization: Bearer <token_admin>"
```

---

#### `DELETE /antrian/:id` — Admin
Hapus atau batalkan antrian. Hanya bisa menghapus antrian di cabang sendiri.

---

#### `GET /cabang/:id/antrian/detail` — Admin
Lihat **data lengkap** semua antrian di cabang admin yang login (termasuk No HP, No Rangka, dll).

**🧪 cURL:**
```bash
curl "http://localhost:8080/api/cabang/1/antrian/detail?status=menunggu" \
  -H "Authorization: Bearer <token_admin>"
```

---

### USER MANAGEMENT

---

#### `GET /admin/users/kontak` — Admin
Ambil daftar nama dan nomor WhatsApp dari seluruh user (role `user`). Digunakan untuk list kontak atau broadcast manual.

**🧪 cURL:**
```bash
curl http://localhost:8080/api/admin/users/kontak \
  -H "Authorization: Bearer <token_admin>"
```

**📤 Response `200`:**
```json
{
  "success": true,
  "data": [
    {
      "id": 2,
      "name": "Budi Google",
      "no_wa": "081234567890"
    }
  ]
}
```

---

#### `GET /users` — Admin
Ambil daftar semua user yang terdaftar di sistem.

**🧪 cURL:**
```bash
curl http://localhost:8080/api/users \
  -H "Authorization: Bearer <token_admin>"
```

**📤 Response `200`:**
```json
{
  "success": true,
  "data": [
    { "id": 1, "name": "Admin Antrian", "username": "admin_antrian", "role": "admin", "cabang_id": null },
    { "id": 2, "name": "Budi Google", "email": "budi@gmail.com", "role": "user", "cabang_id": null }
  ]
}
```

---

#### `POST /users` — Admin
Buat user atau admin baru. Password akan **otomatis di-hash (bcrypt)** oleh server sebelum disimpan.

**📥 Request Body:**
```json
{
  "name": "Admin Kedaton",
  "username": "admin_kedaton",
  "password": "Admin@Kedaton123",
  "role": "admin",
  "cabang_id": 1
}
```

| Field | Tipe | Required | Keterangan |
|---|---|---|---|
| `name` | string | ✅ | Nama lengkap |
| `username` | string | ⚠️ | Wajib jika `role: "admin"` — digunakan untuk login |
| `email` | string | ⚠️ | Wajib jika `role: "user"` — untuk Google login |
| `password` | string | ✅ | Min. 6 karakter. Disimpan sebagai bcrypt hash |
| `role` | string | ✅ | `"user"` atau `"admin"` |
| `cabang_id` | int | ⚠️ | Wajib jika `role: "admin"` — ID cabang yang dikelola |

**📤 Response `201`:**
```json
{
  "id": 5,
  "name": "Admin Kedaton",
  "username": "admin_kedaton",
  "role": "admin",
  "cabang_id": 1,
  "created_at": "2026-05-01T13:53:57Z"
}
```

**🧪 cURL:**
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token_admin>" \
  -d '{
    "name": "Admin Kedaton",
    "username": "admin_kedaton",
    "password": "Admin@Kedaton123",
    "role": "admin",
    "cabang_id": 1
  }'
```

> ⚠️ Endpoint ini membutuhkan token admin yang sudah login. Untuk membuat admin pertama, gunakan SQL Seeder (lihat Seksi 12).

---

#### `PUT /users/:id` — Admin
Update data user. Jika field `password` diisi, password lama akan diganti (otomatis di-hash ulang).

---

#### `DELETE /users/:id` — Admin
Hapus user dari sistem.

---

### MANAJEMEN SUPER ADMIN (Khusus Global Admin)

Grup endpoint `/api/super/*` diperuntukkan khusus bagi Super Admin. **Super Admin** adalah *user* dengan peran (`role`) sebagai `"admin"` tetapi **tidak terikat pada cabang mana pun** (`cabang_id` bernilai `null`).

---

#### `POST /super/admin/cabang` — Super Admin
Membuat akun admin cabang yang baru.

**📥 Request Body:**
```json
{
  "name": "Admin Pahoman",
  "username": "admin_pahoman",
  "password": "Password123",
  "cabang_id": 2
}
```
*Catatan: Parameter `role` akan otomatis di-_set_ menjadi `"admin"` di sisi backend.*

---

#### `GET /super/admins` — Super Admin
Lihat daftar keseluruhan semua admin (baik Super Admin maupun admin cabang biasa) yang ada di sistem.

---

#### `GET /super/cabang/:id/admins` — Super Admin
Melihat daftar akun admin yang ditugaskan khusus di ID cabang tersebut.

---

#### `PUT /super/admins/:id/assign` — Super Admin
Menugaskan (*assign*) ulang sebuah akun admin ke cabang lain.

**📥 Request Body:**
```json
{
  "cabang_id": 3
}
```

---

#### `DELETE /super/admins/:id/assign` — Super Admin
Melepas (*un-assign*) tugas akun admin dari cabangnya. Kolom `cabang_id` pengguna tersebut akan diubah menjadi `null`, yang secara teknis dapat menjadikannya Super Admin.

---

### KONTAK WA USER

---

#### `GET /users/kontak` — User (wajib login)
Ambil nomor WhatsApp yang tersimpan untuk user yang sedang login.

**🧪 cURL:**
```bash
curl http://localhost:8080/api/users/kontak \
  -H "Authorization: Bearer <token_user>"
```

**📤 Response `200`:**
```json
{
  "success": true,
  "data": {
    "user_id": 2,
    "no_wa": "08123456789"
  }
}
```

> Jika nomor belum pernah disimpan, `no_wa` akan bernilai string kosong `""`.

---

#### `POST /users/kontak` — User (wajib login)
Simpan atau update nomor WhatsApp user yang sedang login. Bisa juga menggunakan `PUT /users/kontak`.

**📥 Request Body:**
```json
{
  "no_wa": "08123456789"
}
```

| Field | Tipe | Required | Keterangan |
|---|---|---|---|
| `no_wa` | string | ✅ | Nomor WA user (format bebas, misal: `08xx`, `62xx`) |

**🧪 cURL:**
```bash
curl -X POST http://localhost:8080/api/users/kontak \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token_user>" \
  -d '{"no_wa": "08123456789"}'
```

**📤 Response `200`:**
```json
{
  "success": true,
  "message": "Kontak WA berhasil disimpan",
  "data": {
    "user_id": 2,
    "no_wa": "08123456789"
  }
}
```

**📤 Response `400` (no_wa kosong):**
```json
{ "success": false, "message": "no_wa wajib diisi" }
```

---

#### `PUT /users/kontak` — User (wajib login)
Alias dari `POST /users/kontak`. Fungsi identik — simpan atau perbarui nomor WA.

**🧪 cURL:**
```bash
curl -X PUT http://localhost:8080/api/users/kontak \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token_user>" \
  -d '{"no_wa": "08199999999"}'
```

---

#### `DELETE /users/kontak` — User (wajib login)
Hapus (kosongkan) nomor WhatsApp user.

**🧪 cURL:**
```bash
curl -X DELETE http://localhost:8080/api/users/kontak \
  -H "Authorization: Bearer <token_user>"
```

**📤 Response `200`:**
```json
{
  "success": true,
  "message": "Kontak WA berhasil dihapus",
  "data": {
    "user_id": 2,
    "no_wa": ""
  }
}
```

---

### CRM — Pengingat WhatsApp

---

#### `POST /crm/send` — Admin
Kirim pesan WA secara manual. Admin bebas tentukan nomor WA dan isi pesan.

**📥 Request Body:**
```json
{
  "no_wa": "081234567890",
  "pesan": "Promo servis bulan ini diskon 20%! Segera booking antrian Anda."
}
```

| Field | Tipe | Required |
|---|---|---|
| `no_wa` | string | ✅ |
| `pesan` | string | ✅ |

**🧪 cURL:**
```bash
curl -X POST http://localhost:8080/api/crm/send \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token_admin>" \
  -d '{"no_wa":"081234567890","pesan":"Motor Anda siap diambil!"}'
```

**📤 Response `200`:**
```json
{ "success": true, "message": "Pesan berhasil dikirim ke 081234567890" }
```

---

#### `POST /crm/reminders` — Admin
Kirim pengingat WA ke pelanggan berdasarkan data antrian. No WA diambil otomatis dari field `no_wa_reminder` (atau fallback ke `no_hp` jika kosong).

**📥 Request Body:**
```json
{
  "antrian_id": 5,
  "pesan": "Mohon hadir 30 menit sebelum waktu estimasi servis Anda."
}
```

| Field | Tipe | Required | Keterangan |
|---|---|---|---|
| `antrian_id` | int | ✅ | ID antrian tujuan |
| `pesan` | string | ❌ | Jika kosong, gunakan template default |

**Template default yang dikirim jika `pesan` tidak diisi:**
```
📢 Pengingat dari Bengkel

Halo Budi Santoso!

Kami mengingatkan Anda memiliki antrian servis motor:
🏍️ Honda Vario 150 (2021)
🔢 Nomor Antrian: #5
📅 Tanggal: 20 Nov 2024
⏰ Estimasi Jam: 09:00

Mohon hadir tepat waktu. Terima kasih! 🙏
```

**🧪 cURL:**
```bash
curl -X POST http://localhost:8080/api/crm/reminders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token_admin>" \
  -d '{"antrian_id": 5}'
```

---

## 5. 📊 Status Code

| Kode | Arti | Kapan Terjadi |
|---|---|---|
| `200` | OK | Request berhasil (GET, PUT, DELETE) |
| `201` | Created | Data berhasil dibuat (POST register, ambil antrian) |
| `400` | Bad Request | Field kurang / format salah |
| `401` | Unauthorized | Token tidak ada atau sudah expired |
| `403` | Forbidden | Token valid tapi tidak punya akses |
| `404` | Not Found | Data tidak ditemukan |
| `500` | Internal Server Error | Error di server |

---

## 6. ⚠️ Error Handling

**Format error:**
```json
{ "success": false, "message": "pesan error" }
```

**Contoh error yang mungkin terjadi:**

| Situasi | Status | Pesan |
|---|---|---|
| Field wajib tidak diisi | 400 | `"nama_pemilik is required"` |
| Email sudah terdaftar | 400 | `"email sudah terdaftar"` |
| Email/password salah | 401 | `"email atau password salah"` |
| Token tidak dikirim | 401 | `"Authorization header diperlukan"` |
| Token expired | 401 | `"Token tidak valid atau kadaluarsa"` |
| Akses antrian orang lain | 403 | `"Akses ditolak"` |
| Admin akses cabang lain | 400 | `"antrian bukan milik cabang Anda"` |
| Antrian tidak ditemukan | 404 | `"Antrian tidak ditemukan"` |
| Tidak ada antrian menunggu | 400 | `"tidak ada antrian yang menunggu di cabang ini"` |
| reminder_aktif tanpa no_wa | 500 | `"no_wa_reminder wajib diisi jika reminder_aktif = true"` |
| Fonnte gagal kirim | 500 | `"gagal mengirim WA: ..."` |

---

## 7. 🧱 Skema Data

### Cabang
```json
{
  "id": "int",
  "nama": "string",
  "alamat": "string",
  "kota": "string",
  "no_telp": "string",
  "latitude": "float",
  "longitude": "float",
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

### User
```json
{
  "id": "int",
  "name": "string",
  "email": "string (nullable)",
  "username": "string (nullable, khusus admin)",
  "google_id": "string (nullable, khusus Google login)",
  "avatar_url": "string (nullable) — URL foto profil",
  "alamat": "string (nullable)",
  "kota": "string (nullable)",
  "provinsi": "string (nullable)",
  "kode_pos": "string (nullable)",
  "promo_aktif": "boolean",
  "role": "string (user | admin)",
  "cabang_id": "int (nullable, khusus admin)",
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

> `avatar_url` diambil dari field `picture` pada response Google tokeninfo API. Diperbarui otomatis setiap kali user login jika foto profil Google berubah.

### UserProfileResponse (GET /user/profile)
```json
{
  "id": "int",
  "name": "string",
  "email": "string",
  "avatar_url": "string — URL foto profil",
  "alamat": "string",
  "kota": "string",
  "provinsi": "string",
  "kode_pos": "string",
  "promo_aktif": "boolean",
  "role": "string",
  "created_at": "datetime",
  "antrian": [
    {
      "id": "int",
      "nomor_antrian": "int",
      "status": "string (menunggu | dipanggil | selesai)",
      "tanggal_kedatangan": "datetime",
      "estimasi_jam": "string (HH:MM)",
      "merk_motor": "string",
      "tipe_motor": "string",
      "created_at": "datetime",
      "cabang": {
        "id": "int",
        "nama": "string",
        "alamat": "string",
        "kota": "string",
        "no_telp": "string"
      }
    }
  ]
}
```

### Antrian
```json
{
  "id": "int",
  "cabang_id": "int",
  "user_id": "int (nullable)",
  "nomor_antrian": "int",
  "status": "string (menunggu | dipanggil | selesai)",
  "nama_pemilik": "string",
  "no_hp": "string",
  "merk_motor": "string",
  "tipe_motor": "string",
  "no_rangka": "string",
  "no_mesin": "string",
  "tahun_pembuatan": "int",
  "tanggal_kedatangan": "datetime",
  "estimasi_jam": "string (HH:MM)",
  "reminder_aktif": "bool",
  "no_wa_reminder": "string (nullable)",
  "catatan": "string (nullable)",
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

---

## 8. 🔔 Alur WA Otomatis

| Kejadian | Syarat | Dikirim ke |
|---|---|---|
| Admin `call-next` | `reminder_aktif: true` | `no_wa_reminder` |
| Admin tandai `selesai` | `reminder_aktif: true` | `no_wa_reminder` |
| Admin kirim manual (`/crm/send`) | — | No WA yang diinput admin |
| Admin kirim dari antrian (`/crm/reminders`) | — | `no_wa_reminder` atau fallback `no_hp` |

---

## 9. 🏢 Cabang Lautan Teduh

| # | Nama Cabang | Kota |
|---|---|---|
| 1 | Lautan Teduh Kedaton | Bandar Lampung |
| 2 | Lautan Teduh Pahoman | Bandar Lampung |
| 3 | Lautan Teduh Tirtayasa | Bandar Lampung |
| 4 | Lautan Teduh Pramuka | Bandar Lampung |
| 5 | Lautan Teduh Karang Anyar | Lampung Selatan |
| 6 | Lautan Teduh Purbolinggo | Lampung Timur |
| 7 | Lautan Teduh Pekalongan | Lampung Timur |
| 8 | Lautan Teduh Metro | Metro |
| 9 | Lautan Teduh Kotabumi | Lampung Utara |
| 10 | Lautan Teduh Kalianda | Lampung Selatan |

---

## 10. 📱 Broadcast — In-App Notification

> Fitur ini **berbeda dari WhatsApp**. Notifikasi disimpan di database dan dibaca user melalui inbox di aplikasi Android.

| | WhatsApp (`/crm/*`) | In-App Broadcast (`/broadcast`) |
|---|---|---|
| **Diterima via** | Aplikasi WhatsApp | Inbox di dalam app |
| **Kapan** | Langsung | Saat user buka app / refresh |
| **Bisa bergambar** | ❌ | ✅ |
| **Tipe** | Pesan teks | Promo / per cabang |

---

### `POST /broadcast` — Admin
Kirim broadcast baru. Pilih `tipe: "promo"` untuk semua user, atau `tipe: "antrian"` untuk user di cabang tertentu saja.

**📥 Request Body:**
```json
{
  "judul": "Promo Servis Lebaran 🎉",
  "deskripsi": "Diskon 30% untuk semua jenis servis di bulan Mei",
  "detail": "Dapatkan diskon 30% untuk servis oli, tune-up, dan ganti ban. Berlaku 1–31 Mei 2025 di seluruh cabang Lautan Teduh.",
  "gambar_url": "https://example.com/promo-mei.jpg",
  "tipe": "promo"
}
```

**Untuk broadcast per cabang (antrian):**
```json
{
  "judul": "Info Cabang Kedaton",
  "deskripsi": "Area parkir sedang dalam perbaikan hari ini",
  "detail": "Mohon maaf, area parkir cabang Kedaton sedang dalam perbaikan. Harap gunakan parkir alternatif di sebelah gedung.",
  "gambar_url": "",
  "tipe": "antrian",
  "cabang_id": 1
}
```

| Field | Tipe | Required | Keterangan |
|---|---|---|---|
| `judul` | string | ✅ | Judul notifikasi |
| `deskripsi` | string | ✅ | Teks singkat tampil di list |
| `detail` | string | ✅ | Isi lengkap tampil saat klik |
| `gambar_url` | string | ❌ | URL gambar (kosongkan jika tidak ada) |
| `tipe` | string | ✅ | `"promo"` (semua user) atau `"antrian"` (per cabang) |
| `cabang_id` | int | ⚠️ | Wajib jika `tipe: "antrian"` |

**📤 Response `201`:**
```json
{
  "success": true,
  "message": "Broadcast berhasil dikirim",
  "data": { "id": 3, "judul": "Promo Servis Lebaran 🎉", "tipe": "promo", ... }
}
```

**🧪 cURL:**
```bash
curl -X POST http://localhost:8080/api/broadcast \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token_admin>" \
  -d '{"judul":"Promo Servis","deskripsi":"Diskon 30%","detail":"Detail promo...","tipe":"promo"}'
```

---

### `GET /broadcast` — User (list inbox)
Ambil semua broadcast yang relevan untuk user:
- Semua broadcast `tipe: "promo"` (global)
- Broadcast `tipe: "antrian"` dari cabang yang pernah dipakai user

**🧪 cURL:**
```bash
curl http://localhost:8080/api/broadcast \
  -H "Authorization: Bearer <token_user>"
```

**📤 Response `200` (tampilan list / inbox):**
```json
{
  "success": true,
  "data": [
    {
      "id": 3,
      "judul": "Promo Servis Lebaran 🎉",
      "deskripsi": "Diskon 30% untuk semua jenis servis di bulan Mei",
      "tipe": "promo",
      "cabang_id": null,
      "created_at": "2024-11-20T09:00:00Z"
    },
    {
      "id": 4,
      "judul": "Info Cabang Kedaton",
      "deskripsi": "Area parkir sedang dalam perbaikan hari ini",
      "tipe": "antrian",
      "cabang_id": 1,
      "created_at": "2024-11-19T14:00:00Z"
    }
  ]
}
```

> **Note:** Tampilkan `judul`, `deskripsi`, dan `created_at` di list. **Jangan ambil `detail` dan `gambar_url` di list** untuk menghemat bandwidth.

---

### `GET /broadcast/:id` — User (detail)
Ambil detail lengkap satu broadcast saat user klik notifikasi.

**🧪 cURL:**
```bash
curl http://localhost:8080/api/broadcast/3 \
  -H "Authorization: Bearer <token_user>"
```

**📤 Response `200` (detail setelah diklik):**
```json
{
  "success": true,
  "data": {
    "id": 3,
    "judul": "Promo Servis Lebaran 🎉",
    "deskripsi": "Diskon 30% untuk semua jenis servis di bulan Mei",
    "detail": "Dapatkan diskon 30% untuk servis oli, tune-up, dan ganti ban. Berlaku 1–31 Mei 2025 di seluruh cabang Lautan Teduh.",
    "gambar_url": "https://example.com/promo-mei.jpg",
    "tipe": "promo",
    "cabang_id": null,
    "created_at": "2024-11-20T09:00:00Z"
  }
}
```

> Jika `gambar_url` kosong `""`, tidak perlu tampilkan komponen gambar di UI Android.

---

### `GET /broadcast/all` — Admin
Lihat semua broadcast yang pernah dikirim (semua tipe, semua cabang).

---

### Skema Broadcast
```json
{
  "id": "int",
  "admin_id": "int",
  "judul": "string",
  "deskripsi": "string",
  "detail": "string",
  "gambar_url": "string (nullable)",
  "tipe": "string (promo | antrian)",
  "cabang_id": "int (nullable — null jika promo global)",
  "created_at": "datetime"
}
```

---

## 11. 🗄️ Migrasi Database

Migrasi dilakukan secara manual menggunakan file SQL di folder `migrations/`.

| File | Keterangan |
|---|---|
| `001_init.sql` | Schema awal: tabel `cabangs`, `users`, `antrians`, `customers` + index |
| `002_add_avatar_url.sql` | Tambah kolom `avatar_url TEXT` ke tabel `users` |

### Cara Menjalankan Migrasi (via Docker)

```bash
# Migrasi pertama kali (schema awal)
docker exec -i database_ta_lautan_teduh psql -U admin -d ta_lautan_db \
  < migrations/001_init.sql

# Migrasi tambah kolom avatar_url (jalankan setelah update terbaru)
docker exec -i database_ta_lautan_teduh psql -U admin -d ta_lautan_db \
  < migrations/002_add_avatar_url.sql
```

> Semua migrasi menggunakan `IF NOT EXISTS` / `ADD COLUMN IF NOT EXISTS` sehingga **aman dijalankan berulang kali** tanpa error.

### Verifikasi Kolom avatar_url
```bash
docker exec -it database_ta_lautan_teduh psql -U admin -d ta_lautan_db \
  -c "\d users"
```
Pastikan ada baris: `avatar_url | text`

---

## 12. 🚀 Deployment (VPS + Docker)

Aplikasi ini siap dijalankan di VPS menggunakan **Docker Compose** dengan stack:
- **App** — Go binary (Alpine)
- **Database** — PostgreSQL 15
- **Reverse Proxy** — Nginx (dengan SSL Let's Encrypt)

---

### Struktur File Deployment

```
deployments/
├── Dockerfile           # Build Go binary
└── docker-compose.yml   # Orchestrasi semua service

nginx.conf               # Konfigurasi reverse proxy
.env.example             # Template environment variables
.env                     # File env asli (TIDAK di-commit ke git)
```

---

### Environment Variables

Buat file `.env` di root project (salin dari `.env.example`):

```env
# Database
DB_HOST=db
DB_PORT=5432
DB_USER=admin
DB_PASS=ta-lautan
DB_NAME=ta_lautan_db

# JWT Secret — gunakan string acak yang panjang!
JWT_SECRET=ganti-dengan-secret-yang-sangat-panjang-dan-random

# Fonnte API (WhatsApp)
FONNTE_TOKEN=token-fonnte-anda-disini
```

> ⚠️ **PENTING:** File `.env` sudah masuk `.gitignore`. Jangan pernah commit file `.env` ke repository!

---

### Cara Deploy di VPS

**1. Clone repository ke VPS:**
```bash
git clone https://github.com/rakaascode/server-antrian-go.git
cd server-antrian-go
```

**2. Buat file `.env`:**
```bash
cp .env.example .env
nano .env
# Isi JWT_SECRET dan FONNTE_TOKEN dengan nilai asli
```

**3. Pastikan SSL certificate sudah ada** (untuk Nginx HTTPS):
```bash
# Install certbot jika belum
sudo apt install certbot
sudo certbot certonly --standalone -d rakaascode.site -d www.rakaascode.site
```

**4. Build & jalankan semua service:**
```bash
cd deployments
docker compose up -d --build
```

**5. Cek status semua container:**
```bash
docker compose ps
```

Output yang diharapkan:
```
NAME                         STATUS
backend_ta_lautan_teduh      running
database_ta_lautan_teduh     running (healthy)
nginx_ta_lautan              running
```

---

### Perintah Docker Berguna

| Perintah | Fungsi |
|---|---|
| `docker compose up -d --build` | Build ulang dan jalankan |
| `docker compose down` | Hentikan semua service |
| `docker compose logs -f app` | Lihat log aplikasi real-time |
| `docker compose logs -f db` | Lihat log database |
| `docker compose restart app` | Restart hanya aplikasi |
| `docker compose exec db psql -U admin ta_lautan_db` | Masuk ke database |

---

### Troubleshooting

| Masalah | Solusi |
|---|---|
| App gagal konek ke DB | Tunggu DB sehat dulu: `docker compose logs db` |
| `JWT_SECRET not set` | Pastikan `.env` sudah dibuat dan diisi |
| `Fonnte error` | Cek `FONNTE_TOKEN` di `.env` sudah benar |
| Nginx 502 Bad Gateway | App belum siap: `docker compose logs app` |
| SSL Certificate error | Jalankan certbot dan pastikan path `/etc/letsencrypt` benar |
| Port 80/443 sudah dipakai | Stop service nginx lain: `sudo systemctl stop nginx` |

---

## 13. 🔐 Pembuatan Admin Pertama (Seeding)

Endpoint `POST /api/users` membutuhkan token admin. Untuk membuat admin pertama saat awal deployment, gunakan SQL langsung ke database:

### Cara: SQL via Docker + pgcrypto

```bash
# Masuk ke container database dan jalankan SQL berikut:
docker exec -it database_ta_lautan_teduh psql -U admin -d ta_lautan_db -c "
CREATE EXTENSION IF NOT EXISTS pgcrypto;
INSERT INTO users (name, username, password, role, created_at, updated_at)
VALUES (
  'Admin Antrian',
  'admin_antrian',
  crypt('Admin@Antrian123', gen_salt('bf')),
  'admin',
  NOW(),
  NOW()
)
ON CONFLICT (username) DO UPDATE SET
  password = crypt('Admin@Antrian123', gen_salt('bf')),
  role = 'admin';
"
```

> `crypt(..., gen_salt('bf'))` menghasilkan hash bcrypt yang kompatibel dengan library Go (`golang.org/x/crypto/bcrypt`). Login akan berfungsi normal.

**Verifikasi setelah seeding:**
```bash
curl -X POST http://localhost:8080/api/auth/admin/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin_antrian","password":"Admin@Antrian123"}'
```

Response berhasil:
```json
{
  "success": true,
  "message": "Login admin berhasil",
  "data": {
    "token": "eyJhbGci...",
    "user": { "id": 1, "name": "Admin Antrian", "username": "admin_antrian", "role": "admin" }
  }
}
```

**Setelah punya token**, buat admin untuk cabang lain via API:
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token_dari_login>" \
  -d '{
    "name": "Admin Kedaton",
    "username": "admin_kedaton",
    "password": "Admin@Kedaton123",
    "role": "admin",
    "cabang_id": 1
  }'
```

### Daftar Admin per Cabang (Contoh)

| Cabang | Username | Password Default |
|---|---|---|
| Kedaton | `admin_kedaton` | `Admin@Kedaton123` |
| Pahoman | `admin_pahoman` | `Admin@Pahoman123` |
| Tirtayasa | `admin_tirtayasa` | `Admin@Tirtayasa123` |
| Pramuka | `admin_pramuka` | `Admin@Pramuka123` |
| Karang Anyar | `admin_karanganyar` | `Admin@KarangAnyar123` |

> ⚠️ **Ganti password default** segera setelah pertama kali login!

---

### Update Aplikasi

```bash
# Pull perubahan terbaru
git pull origin main

# Build ulang image dan restart
cd deployments
docker compose up -d --build app

# Verifikasi
docker compose logs -f app
```
