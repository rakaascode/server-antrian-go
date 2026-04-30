package crm

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
