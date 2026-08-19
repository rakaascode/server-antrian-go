# 🔄 Panduan Update & Maintenance Project di VPS

Dokumen ini berisi panduan praktis untuk melakukan update kode, deploy ulang, dan maintenance project secara independen di VPS tanpa mematikan atau mengganggu project lain.

---

## 📌 Prinsip Kerja

Karena setiap project memiliki folder dan `docker-compose.yml` terpisah:
- Update di **Project A** tidak akan me-restart **Project B**.
- Gateway Nginx (`gateway_nginx`) tetap hidup melayani traffic web.
- Database (`db`) tidak perlu di-restart jika hanya memperbarui kode backend.

---

## 🚀 1. Cara Update Project Backend (`server-antrian-go`)

Jika Anda telah melakukan `git push` dari laptop dan ingin memperbarui kode di server VPS:

### Langkah di VPS:
```bash
# 1. Masuk ke direktori repository
cd /home/casper/server-antrian-go

# 2. Tarik kode terbaru dari GitHub
git pull origin main

# 3. Masuk ke folder deployments & rebuild HANYA container app (backend)
cd deployments
docker compose up -d --build app
```

> 💡 **Kenapa hanya `app`?**
> Menjalankan `docker compose up -d --build app` hanya me-rebuild backend Go, membiarkan container database Postgres (`db`) tetap aktif dan sehat.

---

## 🌐 2. Cara Update Project Frontend / Project Lain

Misal Anda memiliki project frontend React / Next.js / Vue di folder `/home/casper/frontend-app`:

```bash
# 1. Masuk ke folder project
cd /home/casper/frontend-app

# 2. Tarik update terbaru
git pull origin main

# 3. Rebuild dan restart container frontend
docker compose up -d --build
```

---

## ⚡ 3. Cara Cepat: Update 1 Perintah dari Laptop (Tanpa Login VPS)

Anda bisa menjalankan update langsung dari terminal lokal di laptop tanpa perlu masuk SSH manual:

```bash
# Update server antrian langsung dari laptop:
ssh casper@103.253.213.216 "cd /home/casper/server-antrian-go && git pull origin main && cd deployments && docker compose up -d --build app"
```

---

## 🔍 4. Cek Status & Log Setelah Update

### Melihat status semua container:
```bash
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
```

### Melihat log container aplikasi (cek error runtime / panic):
```bash
docker logs -f backend_ta_lautan_teduh --tail 50
```

### Melihat log Gateway Nginx:
```bash
docker logs -f gateway_nginx --tail 50
```

---

## 🛠️ 5. Troubleshooting Cepat

| Masalah | Penyebab | Solusi |
| :--- | :--- | :--- |
| Perubahan kode belum muncul | Docker masih pakai image cache lama | Jalankan: `docker compose build --no-cache app && docker compose up -d app` |
| Container status `Restarting` | Ada error fatal / panic pada kode baru | Cek log: `docker logs backend_ta_lautan_teduh` |
| Error `port already allocated` | Ada container lama yang belum mati | Matikan container lama: `docker stop <nama_container>` |
