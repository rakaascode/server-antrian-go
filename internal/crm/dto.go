package crm

import "time"

// ManualSendRequest untuk kirim WA manual (admin input no WA + pesan)
type ManualSendRequest struct {
	NoWA  string `json:"no_wa" binding:"required"`
	Pesan string `json:"pesan" binding:"required"`
}

// ReminderFromAntrianRequest untuk kirim reminder ke antrian tertentu
type ReminderFromAntrianRequest struct {
	AntrianID uint   `json:"antrian_id" binding:"required"`
	Pesan     string `json:"pesan"` // opsional, jika kosong pakai template default
}

// AntrianCrmItem info antrian yang ditampilkan ke admin untuk keperluan CRM
type AntrianCrmItem struct {
	ID            uint      `json:"id"`
	NomorAntrian  int       `json:"nomor_antrian"`
	Status        string    `json:"status"`
	NamaPemilik   string    `json:"nama_pemilik"`
	NoPolisi      string    `json:"no_polisi"`
	NoWAReminder  string    `json:"no_wa_reminder,omitempty"`
	ReminderAktif bool      `json:"reminder_aktif"`
	MerkMotor     string    `json:"merk_motor"`
	TipeMotor     string    `json:"tipe_motor"`
	EstimasiJam   string    `json:"estimasi_jam"`
	Tanggal       time.Time `json:"tanggal_kedatangan"`
}
