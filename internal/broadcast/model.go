package broadcast

import "time"

const (
	TipePromo   = "promo"   // broadcast global ke semua user
	TipeAntrian = "antrian" // broadcast hanya ke user di cabang tertentu
)

type Broadcast struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	AdminID   uint   `json:"admin_id"` // siapa yang mengirim
	Judul     string `json:"judul"`
	Deskripsi string `json:"deskripsi"` // teks pendek untuk preview list
	Detail    string `json:"detail"`    // isi lengkap (ditampilkan saat klik)
	GambarURL string `json:"gambar_url,omitempty"`

	// "promo" = semua user, "antrian" = hanya user di cabang ini
	Tipe     string `json:"tipe"`
	CabangID *uint  `json:"cabang_id,omitempty"` // null jika promo global

	CreatedAt time.Time `json:"created_at"`
}
