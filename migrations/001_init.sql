-- migrations/001_init.sql
-- Full schema with multi-cabang support

CREATE TABLE IF NOT EXISTS cabangs (
    id         SERIAL PRIMARY KEY,
    nama       VARCHAR(150) NOT NULL,
    alamat     TEXT,
    kota       VARCHAR(100),
    no_telp    VARCHAR(20),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS users (
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(100) NOT NULL,
    username   VARCHAR(100) UNIQUE,           -- khusus admin
    email      VARCHAR(150) UNIQUE,           -- user biasa
    google_id  VARCHAR(255) UNIQUE,           -- user Android
    password   VARCHAR(255),
    role       VARCHAR(20) DEFAULT 'user',    -- user | admin
    cabang_id  INT REFERENCES cabangs(id),    -- hanya admin
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS antrians (
    id                   SERIAL PRIMARY KEY,
    cabang_id            INT NOT NULL REFERENCES cabangs(id),
    user_id              INT REFERENCES users(id),   -- nullable
    nomor_antrian        INT NOT NULL,
    status               VARCHAR(20) DEFAULT 'menunggu',  -- menunggu|dipanggil|selesai

    -- Identitas pemilik (sesuai STNK)
    nama_pemilik         VARCHAR(150) NOT NULL,

    -- Data kendaraan
    merk_motor           VARCHAR(100) NOT NULL,
    tipe_motor           VARCHAR(100) NOT NULL,
    no_rangka            VARCHAR(100) NOT NULL,
    no_mesin             VARCHAR(100) NOT NULL,
    tahun_pembuatan      INT NOT NULL,

    -- Jadwal
    tanggal_kedatangan   DATE NOT NULL,
    estimasi_jam         VARCHAR(10) NOT NULL,         -- e.g. "09:00"

    catatan              TEXT,
    created_at           TIMESTAMPTZ DEFAULT NOW(),
    updated_at           TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS customers (
    id         SERIAL PRIMARY KEY,
    nama       VARCHAR(100) NOT NULL,
    no_hp      VARCHAR(20) UNIQUE NOT NULL,
    email      VARCHAR(150),
    alamat     TEXT,
    catatan    TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Index untuk query cepat
CREATE INDEX IF NOT EXISTS idx_antrian_cabang ON antrians(cabang_id);
CREATE INDEX IF NOT EXISTS idx_antrian_user ON antrians(user_id);
CREATE INDEX IF NOT EXISTS idx_antrian_status ON antrians(status);
