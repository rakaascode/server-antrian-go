-- migrations/003_add_no_wa.sql
-- Tambah kolom no_wa (nomor WhatsApp) ke tabel users

ALTER TABLE users ADD COLUMN IF NOT EXISTS no_wa VARCHAR(20) DEFAULT '';
