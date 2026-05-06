package user

import "time"

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Username  string    `json:"username,omitempty" gorm:"uniqueIndex"`  // untuk admin login
	Email     string    `json:"email,omitempty" gorm:"uniqueIndex"`
	GoogleID  string    `json:"google_id,omitempty" gorm:"uniqueIndex"` // untuk user Android
	AvatarURL string    `json:"avatar_url"`                             // foto profil dari Google
	Password  string    `json:"-"`
	Role      string    `json:"role" gorm:"default:'user'"` // user | admin
	CabangID  *uint     `json:"cabang_id,omitempty"`        // hanya admin yang punya
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
