package crm

import "time"

type Customer struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Nama      string    `json:"nama"`
	NoHP      string    `json:"no_hp" gorm:"uniqueIndex"`
	Email     string    `json:"email"`
	Alamat    string    `json:"alamat"`
	Catatan   string    `json:"catatan"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
