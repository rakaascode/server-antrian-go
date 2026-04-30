package cabang

import "time"

type Cabang struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Nama      string    `json:"nama"`
	Alamat    string    `json:"alamat"`
	Kota      string    `json:"kota"`
	NoTelp    string    `json:"no_telp"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
