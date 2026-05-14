package user

import (
	"github.com/rakaascode/server-antrian-go.git/internal/cabang"
	"github.com/rakaascode/server-antrian-go.git/pkg/utils"
)

type UserService interface {
	GetAll() ([]User, error)
	GetByID(id uint) (User, error)
	GetProfile(id uint) (UserProfileResponse, error)
	UpdateProfile(userID uint, req UpdateProfileRequest) (User, error)
	Create(user User) (User, error)
	Update(id uint, user User) (User, error)
	Delete(id uint) error
	SaveKontak(userID uint, noWA string) (User, error)
	DeleteKontak(userID uint) (User, error)
	// Manajemen admin cabang
	CreateAdminCabang(req CreateAdminCabangRequest) (User, error)
	GetAllAdmins() ([]User, error)
	GetAdminsByCabang(cabangID uint) ([]User, error)
	AssignCabang(adminID uint, cabangID uint) (User, error)
	UnassignCabang(adminID uint) (User, error)
	GetAllUsersKontak() ([]User, error)
}

type userService struct {
	repo UserRepository
}

func NewUserService(r UserRepository) UserService {
	return &userService{r}
}

func (s *userService) GetAll() ([]User, error) {
	return s.repo.FindAll()
}

func (s *userService) GetByID(id uint) (User, error) {
	return s.repo.FindByID(id)
}

// GetProfile mengembalikan profil user (nama, email, avatar)
// beserta riwayat antrian yang pernah diambil lengkap dengan info cabang
func (s *userService) GetProfile(id uint) (UserProfileResponse, error) {
	u, err := s.repo.FindByID(id)
	if err != nil {
		return UserProfileResponse{}, err
	}
	u.Password = ""

	antrians, err := s.repo.FindAntrianByUserID(id)
	if err != nil {
		antrians = []AntrianWithCabang{}
	}

	return UserProfileResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		AvatarURL: u.AvatarURL,
		Alamat:    u.Alamat,
		Kota:      u.Kota,
		Provinsi:  u.Provinsi,
		KodePos:   u.KodePos,
		PromoAktif: u.PromoAktif,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
		Antrian:   antrians,
	}, nil
}

func (s *userService) Create(u User) (User, error) {
	// Hash password sebelum disimpan ke database
	if u.Password != "" {
		hashed, err := utils.HashPassword(u.Password)
		if err != nil {
			return User{}, err
		}
		u.Password = hashed
	}
	return s.repo.Create(u)
}

func (s *userService) Update(id uint, u User) (User, error) {
	u.ID = id
	// Hash password baru jika diisi
	if u.Password != "" {
		hashed, err := utils.HashPassword(u.Password)
		if err != nil {
			return User{}, err
		}
		u.Password = hashed
	}
	return s.repo.Update(u)
}

// UpdateProfile khusus untuk user memperbarui profilnya (alamat, avatar, nama)
func (s *userService) UpdateProfile(userID uint, req UpdateProfileRequest) (User, error) {
	u, err := s.repo.FindByID(userID)
	if err != nil {
		return User{}, err
	}

	// Update field yang diperbolehkan
	if req.Name != "" {
		u.Name = req.Name
	}
	if req.AvatarURL != "" {
		u.AvatarURL = req.AvatarURL
	}
	u.Alamat = req.Alamat
	u.Kota = req.Kota
	u.Provinsi = req.Provinsi
	u.KodePos = req.KodePos
	if req.PromoAktif != nil {
		u.PromoAktif = *req.PromoAktif
	}

	return s.repo.Update(u)
}

func (s *userService) Delete(id uint) error {
	return s.repo.Delete(id)
}

// SaveKontak simpan/update nomor WA user
func (s *userService) SaveKontak(userID uint, noWA string) (User, error) {
	return s.repo.SaveKontak(userID, noWA)
}

// DeleteKontak hapus nomor WA user
func (s *userService) DeleteKontak(userID uint) (User, error) {
	return s.repo.DeleteKontak(userID)
}

// ──────────────────────────────────────────────────────────────────────────────
// Admin Cabang
// ──────────────────────────────────────────────────────────────────────────────

// CreateAdminCabangRequest payload untuk membuat admin cabang baru
type CreateAdminCabangRequest struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	CabangID *uint  `json:"cabang_id"` // opsional, bisa assign nanti
}

// CreateAdminCabang membuat akun user baru dengan role admin
func (s *userService) CreateAdminCabang(req CreateAdminCabangRequest) (User, error) {
	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		return User{}, err
	}
	newUser := User{
		Name:     req.Name,
		Username: &req.Username,
		Password: hashed,
		Role:     "admin",
		CabangID: req.CabangID,
	}
	return s.repo.Create(newUser)
}

// GetAllAdmins mengembalikan semua user dengan role admin
func (s *userService) GetAllAdmins() ([]User, error) {
	return s.repo.FindAllAdmins()
}

// GetAdminsByCabang mengembalikan admin yang bertugas di cabang tertentu
func (s *userService) GetAdminsByCabang(cabangID uint) ([]User, error) {
	return s.repo.FindAdminsByCabang(cabangID)
}

// AssignCabang menugaskan admin ke sebuah cabang
func (s *userService) AssignCabang(adminID uint, cabangID uint) (User, error) {
	return s.repo.UpdateCabangID(adminID, &cabangID)
}

// UnassignCabang melepas admin dari cabang (cabang_id menjadi NULL)
func (s *userService) UnassignCabang(adminID uint) (User, error) {
	return s.repo.UpdateCabangID(adminID, nil)
}

// AntrianWithCabang data antrian beserta nama cabang
type AntrianWithCabang struct {
	ID                uint          `json:"id"`
	NomorAntrian      int           `json:"nomor_antrian"`
	Status            string        `json:"status"`
	TanggalKedatangan interface{}   `json:"tanggal_kedatangan"`
	EstimasiJam       string        `json:"estimasi_jam"`
	MerkMotor         string        `json:"merk_motor"`
	TipeMotor         string        `json:"tipe_motor"`
	Cabang            *cabang.Cabang `json:"cabang,omitempty"`
	CreatedAt         interface{}   `json:"created_at"`
}

// UserProfileResponse response untuk GET /user/profile
type UserProfileResponse struct {
	ID        uint                `json:"id"`
	Name      string              `json:"name"`
	Email     *string             `json:"email,omitempty"`
	AvatarURL string              `json:"avatar_url"`
	Alamat     *string             `json:"alamat,omitempty"`
	Kota       *string             `json:"kota,omitempty"`
	Provinsi   *string             `json:"provinsi,omitempty"`
	KodePos    *string             `json:"kode_pos,omitempty"`
	PromoAktif bool                `json:"promo_aktif"`
	Role       string              `json:"role"`
	CreatedAt interface{}         `json:"created_at"`
	Antrian   []AntrianWithCabang `json:"antrian"`
}

// GetAllUsersKontak mengambil semua data kontak user
func (s *userService) GetAllUsersKontak() ([]User, error) {
	return s.repo.FindAllUsersKontak()
}
