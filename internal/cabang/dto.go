package cabang

type CreateCabangRequest struct {
	Nama      string  `json:"nama" binding:"required"`
	Alamat    string  `json:"alamat" binding:"required"`
	Kota      string  `json:"kota" binding:"required"`
	NoTelp    string  `json:"no_telp"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type UpdateCabangRequest struct {
	Nama      string  `json:"nama"`
	Alamat    string  `json:"alamat"`
	Kota      string  `json:"kota"`
	NoTelp    string  `json:"no_telp"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// CabangPublicResponse untuk endpoint public
type CabangPublicResponse struct {
	ID            uint    `json:"id"`
	Nama          string  `json:"nama"`
	Kota          string  `json:"kota"`
	NoTelp        string  `json:"no_telp"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	JumlahAntrian int64   `json:"jumlah_antrian"`
}

// AntrianHariIni ringkasan antrian hari ini (embed di response cabang)
type AntrianHariIni struct {
	NomorDipanggil *int   `json:"nomor_dipanggil"` // null jika belum ada yang dipanggil
	EstimasiJam    string `json:"estimasi_jam"`    // estimasi jam antrian yang sedang berjalan
	SisaAntrian    int64  `json:"sisa_antrian"`    // jumlah antrian menunggu hari ini
}

// CabangWithAntrianResponse gabungan data cabang + ringkasan antrian hari ini
type CabangWithAntrianResponse struct {
	ID             uint           `json:"id"`
	Nama           string         `json:"nama"`
	Alamat         string         `json:"alamat"`
	Kota           string         `json:"kota"`
	NoTelp         string         `json:"no_telp"`
	Latitude       float64        `json:"latitude"`
	Longitude      float64        `json:"longitude"`
	AntrianHariIni AntrianHariIni `json:"antrian_hari_ini"`
}
