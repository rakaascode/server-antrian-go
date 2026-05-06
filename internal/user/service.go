package user

import "github.com/rakaascode/server-antrian-go.git/internal/cabang"

type UserService interface {
	GetAll() ([]User, error)
	GetByID(id uint) (User, error)
	GetProfile(id uint) (UserProfileResponse, error)
	Create(user User) (User, error)
	Update(id uint, user User) (User, error)
	Delete(id uint) error
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
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
		Antrian:   antrians,
	}, nil
}

func (s *userService) Create(u User) (User, error) {
	return s.repo.Create(u)
}

func (s *userService) Update(id uint, u User) (User, error) {
	u.ID = id
	return s.repo.Update(u)
}

func (s *userService) Delete(id uint) error {
	return s.repo.Delete(id)
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
	Email     string              `json:"email,omitempty"`
	AvatarURL string              `json:"avatar_url"`
	Role      string              `json:"role"`
	CreatedAt interface{}         `json:"created_at"`
	Antrian   []AntrianWithCabang `json:"antrian"`
}
