-- migrations/005_add_user_address.sql
-- Tambah kolom alamat, kota, provinsi, dan kode pos ke tabel users

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS alamat   TEXT DEFAULT NULL,
    ADD COLUMN IF NOT EXISTS kota     VARCHAR(100) DEFAULT NULL,
    ADD COLUMN IF NOT EXISTS provinsi VARCHAR(100) DEFAULT NULL,
    ADD COLUMN IF NOT EXISTS kode_pos VARCHAR(20) DEFAULT NULL;
