package antrian

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

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
	NomorDisplay string `json:"nomor_display" gorm:"-"` // format "A-001", display only
	Status       string `json:"status" gorm:"default:'menunggu'"`

	// Identitas pemilik kendaraan (sesuai STNK)
	NamaPemilik string `json:"nama_pemilik"`
	NoPolisi    string `json:"no_polisi"` // nomor polisi plat kendaraan

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
	NoWAReminder  string `json:"no_wa_reminder,omitempty"` // nomor WA untuk notifikasi

	Catatan   string    `json:"catatan,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (a *Antrian) formatNomorDisplay() {
	if a.NomorAntrian > 0 {
		a.NomorDisplay = fmt.Sprintf("A-%03d", a.NomorAntrian)
	}
}

func (a *Antrian) AfterFind(tx *gorm.DB) (err error) {
	a.formatNomorDisplay()
	return nil
}

func (a *Antrian) AfterCreate(tx *gorm.DB) (err error) {
	a.formatNomorDisplay()
	return nil
}
