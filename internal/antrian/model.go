package antrian

import "time"

const (
	StatusMenunggu   = "menunggu"
	StatusDipanggil  = "dipanggil"
	StatusSelesai    = "selesai"
	StatusDibatalkan = "dibatalkan"
)

type Antrian struct {
	ID       uint  `json:"id" gorm:"primaryKey"`
	CabangID uint  `json:"cabang_id"`
	UserID   *uint `json:"user_id,omitempty"` // nullable: user yg ambil via app

	NomorAntrian int    `json:"nomor_antrian"`
	Status       string `json:"status" gorm:"default:'menunggu'"`

	// Identitas pemilik kendaraan (sesuai STNK)
	NamaPemilik string `json:"nama_pemilik"`
	NoHP        string `json:"no_hp"` // nomor HP sesuai STNK / kontak utama

	// Data kendaraan
	MerkMotor      string `json:"merk_motor"`
	TipeMotor      string `json:"tipe_motor"`
	NoRangka       string `json:"no_rangka"`
	NoMesin        string `json:"no_mesin"`
	TahunPembuatan int    `json:"tahun_pembuatan"`

	// Jadwal
	TanggalKedatangan time.Time `json:"tanggal_kedatangan"`
	EstimasiJam       string    `json:"estimasi_jam"` // e.g. "09:00"

	// Pengingat WA (opsional — diisi user jika mau dapat notif WA)
	ReminderAktif bool   `json:"reminder_aktif" gorm:"default:false"`
	NoWAReminder  string `json:"no_wa_reminder,omitempty"` // bisa beda dari NoHP

	Catatan   string    `json:"catatan,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
