package user

import "gorm.io/gorm"

type UserRepository interface {
	FindAll() ([]User, error)
	FindByID(id uint) (User, error)
	FindByEmail(email string) (User, error)
	FindByUsername(username string) (User, error)
	FindByGoogleID(googleID string) (User, error)
	Create(user User) (User, error)
	Update(user User) (User, error)
	Delete(id uint) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) FindAll() ([]User, error) {
	var users []User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *userRepository) FindByID(id uint) (User, error) {
	var user User
	err := r.db.First(&user, id).Error
	return user, err
}

func (r *userRepository) FindByEmail(email string) (User, error) {
	var user User
	err := r.db.Where("email = ?", email).First(&user).Error
	return user, err
}

func (r *userRepository) FindByUsername(username string) (User, error) {
	var user User
	err := r.db.Where("username = ?", username).First(&user).Error
	return user, err
}

func (r *userRepository) FindByGoogleID(googleID string) (User, error) {
	var user User
	err := r.db.Where("google_id = ?", googleID).First(&user).Error
	return user, err
}

func (r *userRepository) Create(user User) (User, error) {
	err := r.db.Create(&user).Error
	return user, err
}

func (r *userRepository) Update(user User) (User, error) {
	err := r.db.Save(&user).Error
	return user, err
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&User{}, id).Error
}
