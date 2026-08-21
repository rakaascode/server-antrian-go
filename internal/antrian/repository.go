package antrian

import (
	"gorm.io/gorm"
)

type AntrianRepository interface {
	FindByCabang(cabangID uint) ([]Antrian, error)
	FindByCabangAndStatus(cabangID uint, status string) ([]Antrian, error)
	FindByID(id uint) (Antrian, error)
	FindByUserID(userID uint) ([]Antrian, error)
	Create(a Antrian) (Antrian, error)
	UpdateStatus(id uint, status string) error
	Delete(id uint) error
	// Status antrian realtime
	FindLatestDipanggil(cabangID uint) (*Antrian, error)
	CountMenungguSebelum(cabangID uint, nomorAntrian int) (int64, error)
	// Ringkasan semua cabang (hari ini)
	GetRingkasanSemuaCabang() ([]RingkasanCabangResponse, error)
	// Antrian hari ini yang masih aktif (menunggu/dipanggil) — untuk CRM
	FindTodayActiveByCabang(cabangID uint) ([]Antrian, error)
}

type antrianRepository struct {
	db *gorm.DB
}

func NewAntrianRepository(db *gorm.DB) AntrianRepository {
	return &antrianRepository{db}
}

func (r *antrianRepository) FindByCabang(cabangID uint) ([]Antrian, error) {
	var list []Antrian
	err := r.db.Where("cabang_id = ?", cabangID).Order("nomor_antrian asc").Find(&list).Error
	return list, err
}

func (r *antrianRepository) FindByCabangAndStatus(cabangID uint, status string) ([]Antrian, error) {
	var list []Antrian
	err := r.db.Where("cabang_id = ? AND status = ?", cabangID, status).
		Order("nomor_antrian asc").Find(&list).Error
	return list, err
}

func (r *antrianRepository) FindByID(id uint) (Antrian, error) {
	var a Antrian
	err := r.db.First(&a, id).Error
	return a, err
}

func (r *antrianRepository) FindByUserID(userID uint) ([]Antrian, error) {
	var list []Antrian
	err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&list).Error
	return list, err
}

// Create membuat antrian baru dengan nomor urut yang aman dari race condition.
// Menggunakan pg_advisory_xact_lock per cabang sehingga transaksi menghitung MAX dan INSERT
// secara atomic untuk cabang yang sama.
func (r *antrianRepository) Create(a Antrian) (Antrian, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SELECT pg_advisory_xact_lock(?)", int64(a.CabangID)).Error; err != nil {
			return err
		}

		var maxNomor *int
		if err := tx.Model(&Antrian{}).
			Select("MAX(nomor_antrian)").
			Where("cabang_id = ? AND DATE(created_at) = CURRENT_DATE", a.CabangID).
			Scan(&maxNomor).Error; err != nil {
			return err
		}

		if maxNomor != nil {
			a.NomorAntrian = *maxNomor + 1
		} else {
			a.NomorAntrian = 1
		}

		return tx.Create(&a).Error
	})
	return a, err
}

func (r *antrianRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&Antrian{}).Where("id = ?", id).Update("status", status).Error
}

func (r *antrianRepository) Delete(id uint) error {
	return r.db.Delete(&Antrian{}, id).Error
}

// FindLatestDipanggil ambil antrian yang sedang dipanggil (terbesar nomornya)
func (r *antrianRepository) FindLatestDipanggil(cabangID uint) (*Antrian, error) {
	var a Antrian
	err := r.db.Where("cabang_id = ? AND status = ?", cabangID, StatusDipanggil).
		Order("nomor_antrian desc").
		First(&a).Error
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// CountMenungguSebelum hitung berapa antrian "menunggu" dengan nomor < nomorAntrian (= posisi di antrean)
func (r *antrianRepository) CountMenungguSebelum(cabangID uint, nomorAntrian int) (int64, error) {
	var count int64
	err := r.db.Model(&Antrian{}).
		Where("cabang_id = ? AND status = ? AND nomor_antrian < ?", cabangID, StatusMenunggu, nomorAntrian).
		Count(&count).Error
	return count, err
}

// GetRingkasanSemuaCabang mengambil ringkasan antrian hari ini dari semua cabang
func (r *antrianRepository) GetRingkasanSemuaCabang() ([]RingkasanCabangResponse, error) {
	type row struct {
		CabangID       uint
		NamaCabang     string
		Latitude       float64
		Longitude      float64
		NomorDipanggil *int
		EstimasiJam    string
		SisaAntrian    int64
	}

	var rows []row
	err := r.db.Raw(`
		SELECT
			c.id                                        AS cabang_id,
			c.nama                                      AS nama_cabang,
			c.latitude                                  AS latitude,
			c.longitude                                 AS longitude,
			dp.nomor_antrian                            AS nomor_dipanggil,
			dp.estimasi_jam                             AS estimasi_jam,
			COALESCE(wt.total, 0)                       AS sisa_antrian
		FROM cabangs c
		LEFT JOIN LATERAL (
			SELECT nomor_antrian, estimasi_jam
			FROM antrians
			WHERE cabang_id = c.id
			  AND status = 'dipanggil'
			ORDER BY nomor_antrian DESC
			LIMIT 1
		) dp ON true
		LEFT JOIN (
			SELECT cabang_id, COUNT(*) AS total
			FROM antrians
			WHERE status = 'menunggu'
			  AND DATE(tanggal_kedatangan) = CURRENT_DATE
			GROUP BY cabang_id
		) wt ON wt.cabang_id = c.id
		ORDER BY c.id
	`).Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make([]RingkasanCabangResponse, 0, len(rows))
	for _, r := range rows {
		result = append(result, RingkasanCabangResponse{
			CabangID:       r.CabangID,
			NamaCabang:     r.NamaCabang,
			Latitude:       r.Latitude,
			Longitude:      r.Longitude,
			NomorDipanggil: r.NomorDipanggil,
			EstimasiJam:    r.EstimasiJam,
			SisaAntrian:    r.SisaAntrian,
		})
	}
	return result, nil
}

// FindTodayActiveByCabang ambil antrian hari ini yang berstatus menunggu atau dipanggil
func (r *antrianRepository) FindTodayActiveByCabang(cabangID uint) ([]Antrian, error) {
	var list []Antrian
	err := r.db.Where(
		"cabang_id = ? AND status IN ? AND DATE(tanggal_kedatangan) = CURRENT_DATE",
		cabangID,
		[]string{StatusMenunggu, StatusDipanggil},
	).Order("nomor_antrian asc").Find(&list).Error
	return list, err
}
