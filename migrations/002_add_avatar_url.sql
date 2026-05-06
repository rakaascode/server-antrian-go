-- migrations/002_add_avatar_url.sql
-- Tambah kolom avatar_url ke tabel users untuk menyimpan foto profil dari Google

ALTER TABLE users ADD COLUMN IF NOT EXISTS avatar_url TEXT;
