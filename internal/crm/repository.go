package crm

import "gorm.io/gorm"

type CustomerRepository interface {
	FindAll() ([]Customer, error)
	FindByID(id uint) (Customer, error)
	FindByNoHP(noHP string) (Customer, error)
	Create(c Customer) (Customer, error)
	Update(c Customer) (Customer, error)
	Delete(id uint) error
}

type customerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepository{db}
}

func (r *customerRepository) FindAll() ([]Customer, error) {
	var list []Customer
	err := r.db.Order("created_at desc").Find(&list).Error
	return list, err
}

func (r *customerRepository) FindByID(id uint) (Customer, error) {
	var c Customer
	err := r.db.First(&c, id).Error
	return c, err
}

func (r *customerRepository) FindByNoHP(noHP string) (Customer, error) {
	var c Customer
	err := r.db.Where("no_hp = ?", noHP).First(&c).Error
	return c, err
}

func (r *customerRepository) Create(c Customer) (Customer, error) {
	err := r.db.Create(&c).Error
	return c, err
}

func (r *customerRepository) Update(c Customer) (Customer, error) {
	err := r.db.Save(&c).Error
	return c, err
}

func (r *customerRepository) Delete(id uint) error {
	return r.db.Delete(&Customer{}, id).Error
}
