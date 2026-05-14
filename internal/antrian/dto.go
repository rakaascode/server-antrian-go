package antrian

import "time"

// AmbilAntrianRequest payload untuk ambil nomor antrian
type AmbilAntrianRequest struct {
	CabangID uint `json:"cabang_id" binding:"required"`

	// Identitas pemilik sesuai STNK
	NamaPemilik string `json:"nama_pemilik" binding:"required"`
	NoHP        string `json:"no_hp" binding:"required"` // nomor HP sesuai STNK

	// Data kendaraan
	MerkMotor      string `json:"merk_motor" binding:"required"`
	TipeMotor      string `json:"tipe_motor" binding:"required"`
	NoRangka       string `json:"no_rangka" binding:"required"`
	NoMesin        string `json:"no_mesin" binding:"required"`
	TahunPembuatan int    `json:"tahun_pembuatan" binding:"required"`

	// Jadwal
	TanggalKedatangan time.Time `json:"tanggal_kedatangan" binding:"required"`
	EstimasiJam       string    `json:"estimasi_jam" binding:"required"` // "09:00"

	// Pengingat WA (opsional)
	ReminderAktif bool   `json:"reminder_aktif"`
	NoWAReminder  string `json:"no_wa_reminder"` // wajib diisi jika reminder_aktif = true

	Catatan string `json:"catatan"`
}

// AntrianPublicResponse untuk endpoint public (tanpa data sensitif)
type AntrianPublicResponse struct {
	ID                uint      `json:"id"`
	NomorAntrian      int       `json:"nomor_antrian"`
	Status            string    `json:"status"`
	EstimasiJam       string    `json:"estimasi_jam"`
	TanggalKedatangan time.Time `json:"tanggal_kedatangan"`
}

// AntrianDetailResponse untuk pemilik dan admin (data lengkap)
type AntrianDetailResponse struct {
	ID       uint  `json:"id"`
	CabangID uint  `json:"cabang_id"`
	UserID   *uint `json:"user_id,omitempty"`

	NomorAntrian int    `json:"nomor_antrian"`
	Status       string `json:"status"`

	NamaPemilik    string `json:"nama_pemilik"`
	NoHP           string `json:"no_hp"`
	MerkMotor      string `json:"merk_motor"`
	TipeMotor      string `json:"tipe_motor"`
	NoRangka       string `json:"no_rangka"`
	NoMesin        string `json:"no_mesin"`
	TahunPembuatan int    `json:"tahun_pembuatan"`

	TanggalKedatangan time.Time `json:"tanggal_kedatangan"`
	EstimasiJam       string    `json:"estimasi_jam"`

	ReminderAktif bool   `json:"reminder_aktif"`
	NoWAReminder  string `json:"no_wa_reminder,omitempty"`

	Catatan   string    `json:"catatan,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// StatusCabangResponse info realtime antrian di cabang (public)
type StatusCabangResponse struct {
	NomorDipanggil *int   `json:"nomor_dipanggil"` // null jika belum ada yang dipanggil
	StatusPanggil  string `json:"status_panggil"`  // "dipanggil" atau "belum ada"
	TotalMenunggu  int64  `json:"total_menunggu"`
}

// PosisiResponse posisi user di antrian
type PosisiResponse struct {
	NomorAntrian   int    `json:"nomor_antrian"`
	Status         string `json:"status"`
	Posisi         int64  `json:"posisi"`          // berapa orang di depan
	NomorDipanggil *int   `json:"nomor_dipanggil"` // nomor yang sedang dilayani
	Pesan          string `json:"pesan"`
}

// RingkasanCabangResponse ringkasan antrian hari ini per cabang (public)
type RingkasanCabangResponse struct {
	CabangID       uint    `json:"cabang_id"`
	NamaCabang     string  `json:"nama_cabang"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	NomorDipanggil *int    `json:"nomor_dipanggil"` // null jika belum ada yang dipanggil
	EstimasiJam    string  `json:"estimasi_jam"`    // estimasi jam antrian yang sedang berjalan
	SisaAntrian    int64   `json:"sisa_antrian"`    // jumlah antrian berstatus "menunggu" hari ini
}
