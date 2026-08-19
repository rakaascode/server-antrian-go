package cabang

import "gorm.io/gorm"

type CabangRepository interface {
	FindAll() ([]Cabang, error)
	FindByID(id uint) (Cabang, error)
	Create(c Cabang) (Cabang, error)
	Update(c Cabang) (Cabang, error)
	Delete(id uint) error
}

type cabangRepository struct {
	db *gorm.DB
}

func NewCabangRepository(db *gorm.DB) CabangRepository {
	return &cabangRepository{db}
}

func (r *cabangRepository) FindAll() ([]Cabang, error) {
	var list []Cabang
	err := r.db.Order("id asc").Find(&list).Error
	return list, err
}

func (r *cabangRepository) FindByID(id uint) (Cabang, error) {
	var c Cabang
	err := r.db.First(&c, id).Error
	return c, err
}

func (r *cabangRepository) Create(c Cabang) (Cabang, error) {
	err := r.db.Create(&c).Error
	return c, err
}

func (r *cabangRepository) Update(c Cabang) (Cabang, error) {
	err := r.db.Save(&c).Error
	return c, err
}

func (r *cabangRepository) Delete(id uint) error {
	return r.db.Delete(&Cabang{}, id).Error
}
