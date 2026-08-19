# 📘 Panduan Manajemen Server VPS (Gateway & Multi-Project)

Dokumen ini berisi panduan langkah-demi-langkah jika Anda ingin **menambahkan project/server baru** di VPS yang sama, atau jika Anda ingin **pindah ke VPS baru** dari nol.

---

## 🏗️ Arsitektur Server Saat Ini

```
[ Internet / Klien ]
        │ (Port 80 / 443 HTTPS)
        ▼
┌────────────────────────────────────────────────────────┐
│  Gateway Nginx Proxy (/home/casper/gateway-proxy)      │
│  - Network: proxy_net (Docker bridge)                  │
│  - Mengelola SSL & Routing Domain/Subdomain            │
└────────┬───────────────────────────────┬───────────────┘
         │ (internal proxy_net)          │ (internal proxy_net)
         ▼                               ▼
┌────────────────────────────────┐ ┌────────────────────────────────┐
│ Project 1: server-antrian-go   │ │ Project 2: Project Baru / FE   │
│ - backend_ta_lautan_teduh:8080 │ │ - container_baru:port_internal │
│ - database_ta_lautan_teduh     │ └────────────────────────────────┘
└────────────────────────────────┘
```

---

## 🚀 Bagian 1: Cara Menambahkan Project / Server Baru di VPS

Misalkan Anda ingin membuat backend baru (Go/Node.js/Python) atau frontend (Next.js/React/Vue) di domain `app.rakaascode.site` atau domain lain `domain-baru.com`.

### Langkah 1: Siapkan Folder Project
Masuk ke VPS dan buat folder untuk project baru:
```bash
mkdir -p /home/casper/my-new-project
cd /home/casper/my-new-project
```

### Langkah 2: Buat `docker-compose.yml` Project Baru
Pastikan project baru Anda terhubung ke network `proxy_net`:

```yaml
services:
  web_app:
    build: .
    container_name: app_project_baru
    restart: always
    # Expose port lokal (opsional) atau biarkan internal
    ports:
      - "127.0.0.1:3000:3000"
    networks:
      - proxy_net

networks:
  proxy_net:
    external: true
```

Jalankan aplikasinya:
```bash
docker compose up -d --build
```

### Langkah 3: Daftarkan Domain di Gateway Nginx
Buat file konfigurasi routing baru di folder gateway:
```bash
nano /home/casper/gateway-proxy/conf.d/project_baru.conf
```

Isi dengan template berikut (ganti domain dan nama container sesuai project):
```nginx
server {
    listen 80;
    server_name app.rakaascode.site;

    client_max_body_size 20M;

    location / {
        proxy_pass         http://app_project_baru:3000;
        proxy_http_version 1.1;

        proxy_set_header Host              $host;
        proxy_set_header X-Real-IP         $remote_addr;
        proxy_set_header X-Forwarded-For   $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # Mendukung Websocket (jika ada)
        proxy_set_header Upgrade           $http_upgrade;
        proxy_set_header Connection        "upgrade";
    }
}
```

### Langkah 4: Reload Gateway
```bash
docker exec gateway_nginx nginx -s reload
```
*Project baru Anda sudah bisa diakses lewat domain/subdomain tanpa mematikan atau mengganggu project lama!*

---

## 🚚 Bagian 2: Panduan Migrasi Pindah ke VPS Baru

Jika Anda membeli VPS baru dan ingin memindahkan semua server ke sana, ikuti 6 langkah mudah berikut:

### 1. Persiapan di VPS Baru (Install Docker)
Jalankan di terminal VPS baru:
```bash
# Update sistem
sudo apt update && sudo apt upgrade -y

# Install Docker & Docker Compose
curl -fsSL https://get.docker.com | sh

# Beri izin user casper menjalankan docker tanpa sudo
sudo usermod -aG docker $USER
newgrp docker
```

### 2. Setup Network Docker Utama
Buat network bridge bersama yang akan dipakai Gateway dan semua project:
```bash
docker network create proxy_net
```

### 3. Buat Folder Gateway Proxy
```bash
mkdir -p /home/casper/gateway-proxy/conf.d
cd /home/casper/gateway-proxy
```

Buat file `/home/casper/gateway-proxy/docker-compose.yml`:
```yaml
services:
  gateway:
    image: nginx:alpine
    container_name: gateway_nginx
    restart: always
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./conf.d:/etc/nginx/conf.d:ro
      - /etc/letsencrypt:/etc/letsencrypt:ro
    networks:
      - proxy_net

networks:
  proxy_net:
    external: true
```

### 4. Clone Project-Project Anda
```bash
cd /home/casper
git clone git@github.com:rakaascode/server-antrian-go.git

# Masuk ke folder deployments & buat file .env
cd /home/casper/server-antrian-go/deployments
nano .env  # (Isi dengan konfigurasi environment production)

# Jalankan project
docker compose up -d --build
```

### 5. Setup SSL (Let's Encrypt) & Konfigurasi Nginx
Arahkan DNS Domain (A Record) di Cloudflare/Provider Domain ke **IP VPS Baru**.

Install Certbot di VPS untuk sertifikat SSL otomatis:
```bash
sudo apt install -y certbot
sudo certbot certonly --standalone -d rakaascode.site -d www.rakaascode.site
```

Salin template konfigurasi Nginx untuk project ke folder `/home/casper/gateway-proxy/conf.d/antrian.conf`:
```nginx
server {
    listen 80;
    server_name rakaascode.site www.rakaascode.site;
    return 301 https://$host$request_uri;
}

server {
    listen 443 ssl;
    server_name rakaascode.site www.rakaascode.site;

    ssl_certificate /etc/letsencrypt/live/rakaascode.site/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/rakaascode.site/privkey.pem;

    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;

    client_max_body_size 10M;

    location / {
        proxy_pass         http://backend_ta_lautan_teduh:8080;
        proxy_http_version 1.1;

        proxy_set_header Host              $host;
        proxy_set_header X-Real-IP         $remote_addr;
        proxy_set_header X-Forwarded-For   $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        proxy_read_timeout 60s;
        proxy_connect_timeout 10s;
    }
}
```

### 6. Jalankan Gateway
```bash
cd /home/casper/gateway-proxy
docker compose up -d
```

*(Opsional)* **Migrasi Database Lama ke Baru:**
Di VPS lama:
```bash
docker exec -t database_ta_lautan_teduh pg_dump -U admin ta_lautan_db > backup.sql
```
Di VPS baru:
```bash
cat backup.sql | docker exec -i database_ta_lautan_teduh psql -U admin -d ta_lautan_db
```

Selesai! Seluruh sistem telah berhasil dipindahkan ke VPS baru dengan rapi.

---

## 🔄 Bagian 3: Cara Update Project Tertentu Tanpa Mengganggu Project Lain

Karena masing-masing project sudah berada di folder dan `docker-compose` terpisah, proses update **100% terisolasi** (tidak perlu mematikan Gateway Nginx atau project lain).

### Contoh 1: Update Project Antrian (`server-antrian-go`)
Cukup masuk ke folder project antrian, pull dari Git, lalu rebuild container app-nya saja:
```bash
# 1. Masuk ke folder repo
cd /home/casper/server-antrian-go

# 2. Tarik kode terbaru dari GitHub
git pull origin main

# 3. Masuk ke folder deployments & rebuild hanya backend-nya
cd deployments
docker compose up -d --build app
```
> 💡 **Tips:** Menggunakan `docker compose up -d --build app` (menentukan service `app`) akan merebuild backend Go tanpa me-restart container `db` (Postgres), sehingga database tidak terputus.

---

### Contoh 2: Update Project Baru / Frontend
Misal Anda ada update di frontend Next.js / React:
```bash
# 1. Masuk ke folder project frontend
cd /home/casper/my-new-project

# 2. Pull kode terbaru
git pull origin main

# 3. Rebuild container frontend
docker compose up -d --build
```

---

### Perintah Cepat Sekali Jalan (One-Liner):
Anda bisa menjalankan update langsung dari laptop via SSH tanpa perlu login manual:
```bash
ssh casper@103.253.213.216 "cd /home/casper/server-antrian-go && git pull origin main && cd deployments && docker compose up -d --build app"
```
