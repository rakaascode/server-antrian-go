package user

import (
	"github.com/rakaascode/server-antrian-go.git/internal/cabang"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindAll() ([]User, error)
	FindByID(id uint) (User, error)
	FindByEmail(email string) (User, error)
	FindByUsername(username string) (User, error)
	FindByGoogleID(googleID string) (User, error)
	FindAntrianByUserID(userID uint) ([]AntrianWithCabang, error)
	Create(user User) (User, error)
	Update(user User) (User, error)
	Delete(id uint) error
	SaveKontak(userID uint, noWA string) (User, error)
	DeleteKontak(userID uint) (User, error)
	// Manajemen admin cabang
	FindAllAdmins() ([]User, error)
	FindAdminsByCabang(cabangID uint) ([]User, error)
	UpdateCabangID(userID uint, cabangID *uint) (User, error)
	FindAllUsersKontak() ([]User, error)
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

// FindAllAdmins mengambil semua user dengan role admin
func (r *userRepository) FindAllAdmins() ([]User, error) {
	var users []User
	err := r.db.Where("role = ?", "admin").Find(&users).Error
	return users, err
}

// FindAdminsByCabang mengambil admin yang terdaftar di cabang tertentu
func (r *userRepository) FindAdminsByCabang(cabangID uint) ([]User, error) {
	var users []User
	err := r.db.Where("role = ? AND cabang_id = ?", "admin", cabangID).Find(&users).Error
	return users, err
}

// UpdateCabangID mengubah cabang_id admin (assign atau unassign)
func (r *userRepository) UpdateCabangID(userID uint, cabangID *uint) (User, error) {
	var u User
	if err := r.db.First(&u, userID).Error; err != nil {
		return User{}, err
	}
	if err := r.db.Model(&u).Update("cabang_id", cabangID).Error; err != nil {
		return User{}, err
	}
	u.CabangID = cabangID
	return u, nil
}

// SaveKontak menyimpan atau memperbarui nomor WA user
func (r *userRepository) SaveKontak(userID uint, noWA string) (User, error) {
	var u User
	if err := r.db.First(&u, userID).Error; err != nil {
		return User{}, err
	}
	u.NoWA = noWA
	if err := r.db.Model(&u).Update("no_wa", noWA).Error; err != nil {
		return User{}, err
	}
	return u, nil
}

// DeleteKontak menghapus (mengosongkan) nomor WA user
func (r *userRepository) DeleteKontak(userID uint) (User, error) {
	var u User
	if err := r.db.First(&u, userID).Error; err != nil {
		return User{}, err
	}
	if err := r.db.Model(&u).Update("no_wa", "").Error; err != nil {
		return User{}, err
	}
	u.NoWA = ""
	return u, nil
}

// FindAntrianByUserID mengambil riwayat antrian user beserta info cabang
func (r *userRepository) FindAntrianByUserID(userID uint) ([]AntrianWithCabang, error) {
	type rawRow struct {
		ID                uint   `gorm:"column:id"`
		NomorAntrian      int    `gorm:"column:nomor_antrian"`
		Status            string `gorm:"column:status"`
		TanggalKedatangan interface{} `gorm:"column:tanggal_kedatangan"`
		EstimasiJam       string `gorm:"column:estimasi_jam"`
		MerkMotor         string `gorm:"column:merk_motor"`
		TipeMotor         string `gorm:"column:tipe_motor"`
		CreatedAt         interface{} `gorm:"column:created_at"`
		// Cabang
		CabangID    uint   `gorm:"column:cabang_id"`
		CabangNama  string `gorm:"column:cabang_nama"`
		CabangAlamat string `gorm:"column:cabang_alamat"`
		CabangKota  string `gorm:"column:cabang_kota"`
		CabangNoTelp string `gorm:"column:cabang_no_telp"`
	}

	var rows []rawRow
	err := r.db.Raw(`
		SELECT
			a.id, a.nomor_antrian, a.status,
			a.tanggal_kedatangan, a.estimasi_jam,
			a.merk_motor, a.tipe_motor, a.created_at,
			c.id    AS cabang_id,
			c.nama  AS cabang_nama,
			c.alamat AS cabang_alamat,
			c.kota  AS cabang_kota,
			c.no_telp AS cabang_no_telp
		FROM antrians a
		LEFT JOIN cabangs c ON c.id = a.cabang_id
		WHERE a.user_id = ?
		ORDER BY a.created_at DESC
	`, userID).Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make([]AntrianWithCabang, 0, len(rows))
	for _, row := range rows {
		item := AntrianWithCabang{
			ID:                row.ID,
			NomorAntrian:      row.NomorAntrian,
			Status:            row.Status,
			TanggalKedatangan: row.TanggalKedatangan,
			EstimasiJam:       row.EstimasiJam,
			MerkMotor:         row.MerkMotor,
			TipeMotor:         row.TipeMotor,
			CreatedAt:         row.CreatedAt,
		}
		if row.CabangID != 0 {
			item.Cabang = &cabang.Cabang{
				ID:     row.CabangID,
				Nama:   row.CabangNama,
				Alamat: row.CabangAlamat,
				Kota:   row.CabangKota,
				NoTelp: row.CabangNoTelp,
			}
		}
		result = append(result, item)
	}
	return result, nil
}

// FindAllUsersKontak mengambil semua name dan no wa users
func (r *userRepository) FindAllUsersKontak() ([]User, error) {
	var users []User
	err := r.db.Select("id", "name", "no_wa").Where("role = ?", "user").Find(&users).Error
	return users, err
}
