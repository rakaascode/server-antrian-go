package cabang

import (
	"log"

	"gorm.io/gorm"
)

var dataCabangAwal = []Cabang{
	{
		Nama:      "Lautan Teduh Kedaton",
		Alamat:    "Jl. Teuku Umar No.15D, Kedaton",
		Kota:      "Bandar Lampung",
		NoTelp:    "081367846069",
		Latitude:  -5.3795,
		Longitude: 105.2610,
	},
	{
		Nama:      "Lautan Teduh Pahoman",
		Alamat:    "Jl. Gatot Subroto No.93, Kedamaian",
		Kota:      "Bandar Lampung",
		NoTelp:    "081379552244",
		Latitude:  -5.4280,
		Longitude: 105.2595,
	},
	{
		Nama:      "Lautan Teduh Tirtayasa",
		Alamat:    "Jl. P. Tirtayasa No.53, Sukabumi",
		Kota:      "Bandar Lampung",
		NoTelp:    "08994300897",
		Latitude:  -5.3972,
		Longitude: 105.2790,
	},
	{
		Nama:      "Lautan Teduh Pramuka",
		Alamat:    "Jl. Pramuka No.17, Rajabasa",
		Kota:      "Bandar Lampung",
		NoTelp:    "081379168327",
		Latitude:  -5.3700,
		Longitude: 105.2405,
	},
	{
		Nama:      "Lautan Teduh Karang Anyar",
		Alamat:    "Karang Sari, Jati Agung",
		Kota:      "Lampung Selatan",
		NoTelp:    "081272977738",
		Latitude:  -5.3005,
		Longitude: 105.2802,
	},
	{
		Nama:      "Lautan Teduh Purbolinggo",
		Alamat:    "Tanjung Inten, Purbolinggo",
		Kota:      "Lampung Timur",
		NoTelp:    "082250092897",
		Latitude:  -5.1002,
		Longitude: 105.6003,
	},
	{
		Nama:      "Lautan Teduh Pekalongan",
		Alamat:    "Jl. AH Nasution No.18, Adirejo",
		Kota:      "Lampung Timur",
		NoTelp:    "0811345786",
		Latitude:  -5.1504,
		Longitude: 105.5201,
	},
	{
		Nama:      "Lautan Teduh Metro",
		Alamat:    "Jl. Jendral Sudirman No.104, Metro Barat",
		Kota:      "Metro",
		NoTelp:    "081271811215",
		Latitude:  -5.1132,
		Longitude: 105.3075,
	},
	{
		Nama:      "Lautan Teduh Kotabumi",
		Alamat:    "Jl. Alamsyah RPN, Kotabumi Selatan",
		Kota:      "Lampung Utara",
		NoTelp:    "082280585558",
		Latitude:  -4.8255,
		Longitude: 104.8762,
	},
	{
		Nama:      "Lautan Teduh Kalianda",
		Alamat:    "Jl. Kolonel Makmun Rasyid No.166",
		Kota:      "Lampung Selatan",
		NoTelp:    "082175011254",
		Latitude:  -5.7250,
		Longitude: 105.5920,
	},
}

// Seed mengisi data cabang awal jika tabel masih kosong
func Seed(db *gorm.DB) {
	var count int64
	db.Model(&Cabang{}).Count(&count)
	if count > 0 {
		log.Printf("ℹ️  Data cabang sudah ada (%d cabang), skip seeding", count)
		return
	}

	if err := db.Create(&dataCabangAwal).Error; err != nil {
		log.Printf("❌ Gagal seed data cabang: %v", err)
		return
	}
	log.Printf("✅ Berhasil seed %d data cabang Lautan Teduh", len(dataCabangAwal))
}
