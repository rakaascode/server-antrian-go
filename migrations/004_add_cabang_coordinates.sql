-- migrations/004_add_cabang_coordinates.sql
-- Tambah kolom latitude & longitude ke tabel cabangs jika belum ada

ALTER TABLE cabangs
    ADD COLUMN IF NOT EXISTS latitude  DOUBLE PRECISION DEFAULT 0,
    ADD COLUMN IF NOT EXISTS longitude DOUBLE PRECISION DEFAULT 0;
