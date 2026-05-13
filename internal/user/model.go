package user

import "time"

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Username  *string   `json:"username,omitempty" gorm:"uniqueIndex"`  // untuk admin login (NULL untuk Google users)
	Email     *string   `json:"email,omitempty" gorm:"uniqueIndex"`
	GoogleID  *string   `json:"google_id,omitempty" gorm:"uniqueIndex"` // untuk user Android (NULL jika bukan Google user)
	AvatarURL string    `json:"avatar_url"`                             // foto profil dari Google
	Alamat    *string   `json:"alamat,omitempty"`
	Kota      *string   `json:"kota,omitempty"`
	Provinsi  *string   `json:"provinsi,omitempty"`
	KodePos   *string   `json:"kode_pos,omitempty"`
	PromoAktif bool     `json:"promo_aktif" gorm:"default:false"` // bersedia dikirim info promo WA
	NoWA      string    `json:"no_wa,omitempty"`                        // nomor WhatsApp user
	Password  string    `json:"-"`
	Role      string    `json:"role" gorm:"default:'user'"` // user | admin
	CabangID  *uint     `json:"cabang_id,omitempty"`        // hanya admin yang punya
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// KontakRequest payload untuk simpan/update nomor WA
type KontakRequest struct {
	NoWA string `json:"no_wa" binding:"required"`
}

// UpdateProfileRequest payload untuk mengupdate profil
type UpdateProfileRequest struct {
	Name      string  `json:"name"`
	AvatarURL string  `json:"avatar_url"`
	Alamat     *string `json:"alamat"`
	Kota       *string `json:"kota"`
	Provinsi   *string `json:"provinsi"`
	KodePos    *string `json:"kode_pos"`
	PromoAktif *bool   `json:"promo_aktif"`
}
