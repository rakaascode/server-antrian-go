package broadcast

import "time"

// CreateBroadcastRequest payload untuk admin kirim broadcast
type CreateBroadcastRequest struct {
	Judul     string `json:"judul" binding:"required"`
	Deskripsi string `json:"deskripsi" binding:"required"` // teks pendek preview
	Detail    string `json:"detail" binding:"required"`    // isi lengkap
	GambarURL string `json:"gambar_url"`                   // opsional

	// "promo" = semua user, "antrian" = per cabang
	Tipe     string `json:"tipe" binding:"required,oneof=promo antrian"`
	CabangID *uint  `json:"cabang_id"` // wajib jika tipe "antrian"
}

// BroadcastListItem tampilan di list (ringkas)
type BroadcastListItem struct {
	ID        uint      `json:"id"`
	Judul     string    `json:"judul"`
	Deskripsi string    `json:"deskripsi"`
	Tipe      string    `json:"tipe"`
	CabangID  *uint     `json:"cabang_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// BroadcastDetail tampilan lengkap saat klik notifikasi
type BroadcastDetail struct {
	ID        uint      `json:"id"`
	AdminID   uint      `json:"admin_id"`
	Judul     string    `json:"judul"`
	Deskripsi string    `json:"deskripsi"`
	Detail    string    `json:"detail"`
	GambarURL string    `json:"gambar_url,omitempty"`
	Tipe      string    `json:"tipe"`
	CabangID  *uint     `json:"cabang_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
