-- migrations/006_add_promo_aktif.sql
-- Tambah flag persetujuan user untuk dikirim pesan promo via WhatsApp

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS promo_aktif BOOLEAN DEFAULT FALSE;
