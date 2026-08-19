package broadcast

import "gorm.io/gorm"

type BroadcastRepository interface {
	Create(b Broadcast) (Broadcast, error)
	FindAll() ([]Broadcast, error)
	// FindForUser: promo global + antrian dari cabang yang diberikan
	FindForUser(cabangIDs []uint) ([]Broadcast, error)
	FindByID(id uint) (Broadcast, error)
	FindByCabang(cabangID uint) ([]Broadcast, error)
	// Delete hapus satu broadcast berdasarkan ID
	Delete(id uint) error
	// DeleteAll hapus seluruh broadcast
	DeleteAll() error
}

type broadcastRepository struct {
	db *gorm.DB
}

func NewBroadcastRepository(db *gorm.DB) BroadcastRepository {
	return &broadcastRepository{db}
}

func (r *broadcastRepository) Create(b Broadcast) (Broadcast, error) {
	err := r.db.Create(&b).Error
	return b, err
}

func (r *broadcastRepository) FindAll() ([]Broadcast, error) {
	var list []Broadcast
	err := r.db.Order("created_at desc").Find(&list).Error
	return list, err
}

// FindForUser ambil semua promo global + antrian khusus cabang user
func (r *broadcastRepository) FindForUser(cabangIDs []uint) ([]Broadcast, error) {
	var list []Broadcast
	query := r.db.Where("tipe = ?", TipePromo)
	if len(cabangIDs) > 0 {
		query = query.Or("tipe = ? AND cabang_id IN ?", TipeAntrian, cabangIDs)
	}
	err := query.Order("created_at desc").Find(&list).Error
	return list, err
}

func (r *broadcastRepository) FindByID(id uint) (Broadcast, error) {
	var b Broadcast
	err := r.db.First(&b, id).Error
	return b, err
}

func (r *broadcastRepository) FindByCabang(cabangID uint) ([]Broadcast, error) {
	var list []Broadcast
	err := r.db.Where("cabang_id = ? AND tipe = ?", cabangID, TipeAntrian).
		Order("created_at desc").Find(&list).Error
	return list, err
}

// Delete menghapus satu broadcast berdasarkan ID
func (r *broadcastRepository) Delete(id uint) error {
	result := r.db.Delete(&Broadcast{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// DeleteAll menghapus seluruh data broadcast
func (r *broadcastRepository) DeleteAll() error {
	return r.db.Where("1 = 1").Delete(&Broadcast{}).Error
}
