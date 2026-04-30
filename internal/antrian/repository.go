package antrian

import (
	"time"

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
	CountTodayByCabang(cabangID uint) (int64, error)
	// Status antrian realtime
	FindLatestDipanggil(cabangID uint) (*Antrian, error)
	CountMenungguSebelum(cabangID uint, nomorAntrian int) (int64, error)
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

func (r *antrianRepository) Create(a Antrian) (Antrian, error) {
	err := r.db.Create(&a).Error
	return a, err
}

func (r *antrianRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&Antrian{}).Where("id = ?", id).Update("status", status).Error
}

func (r *antrianRepository) Delete(id uint) error {
	return r.db.Delete(&Antrian{}, id).Error
}

func (r *antrianRepository) CountTodayByCabang(cabangID uint) (int64, error) {
	var count int64
	today := time.Now().Truncate(24 * time.Hour)
	err := r.db.Model(&Antrian{}).
		Where("cabang_id = ? AND created_at >= ?", cabangID, today).
		Count(&count).Error
	return count, err
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
