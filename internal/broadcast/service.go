package broadcast

import (
	"errors"

	"github.com/rakaascode/server-antrian-go.git/internal/antrian"
)

type BroadcastService interface {
	// Admin: kirim broadcast
	Create(req CreateBroadcastRequest, adminID uint) (Broadcast, error)
	// Admin: lihat semua broadcast
	GetAll() ([]BroadcastListItem, error)
	// User: ambil broadcast yang relevan (promo + antrian cabang user)
	GetForUser(userID uint) ([]BroadcastListItem, error)
	// Semua: lihat detail satu broadcast
	GetByID(id uint) (BroadcastDetail, error)
	// Admin: hapus satu broadcast by ID
	Delete(id uint) error
	// Admin: hapus semua broadcast
	DeleteAll() error
}

type broadcastService struct {
	repo        BroadcastRepository
	antrianRepo antrian.AntrianRepository
}

func NewBroadcastService(repo BroadcastRepository, antrianRepo antrian.AntrianRepository) BroadcastService {
	return &broadcastService{repo: repo, antrianRepo: antrianRepo}
}

func (s *broadcastService) Create(req CreateBroadcastRequest, adminID uint) (Broadcast, error) {
	// Validasi: tipe antrian wajib isi cabang_id
	if req.Tipe == TipeAntrian && req.CabangID == nil {
		return Broadcast{}, errors.New("cabang_id wajib diisi untuk broadcast tipe antrian")
	}
	// Promo tidak perlu cabang_id
	if req.Tipe == TipePromo {
		req.CabangID = nil
	}

	b := Broadcast{
		AdminID:   adminID,
		Judul:     req.Judul,
		Deskripsi: req.Deskripsi,
		Detail:    req.Detail,
		GambarURL: req.GambarURL,
		Tipe:      req.Tipe,
		CabangID:  req.CabangID,
	}
	return s.repo.Create(b)
}

func (s *broadcastService) GetAll() ([]BroadcastListItem, error) {
	list, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}
	return toListItems(list), nil
}

// GetForUser: ambil broadcast relevan untuk user (promo global + antrian dari cabang user)
func (s *broadcastService) GetForUser(userID uint) ([]BroadcastListItem, error) {
	// Kumpulkan semua cabang_id dari antrian user
	myAntrian, err := s.antrianRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	cabangSet := map[uint]bool{}
	for _, a := range myAntrian {
		cabangSet[a.CabangID] = true
	}
	var cabangIDs []uint
	for id := range cabangSet {
		cabangIDs = append(cabangIDs, id)
	}

	list, err := s.repo.FindForUser(cabangIDs)
	if err != nil {
		return nil, err
	}
	return toListItems(list), nil
}

func (s *broadcastService) GetByID(id uint) (BroadcastDetail, error) {
	b, err := s.repo.FindByID(id)
	if err != nil {
		return BroadcastDetail{}, errors.New("broadcast tidak ditemukan")
	}
	return BroadcastDetail{
		ID:        b.ID,
		AdminID:   b.AdminID,
		Judul:     b.Judul,
		Deskripsi: b.Deskripsi,
		Detail:    b.Detail,
		GambarURL: b.GambarURL,
		Tipe:      b.Tipe,
		CabangID:  b.CabangID,
		CreatedAt: b.CreatedAt,
	}, nil
}

func toListItems(list []Broadcast) []BroadcastListItem {
	result := make([]BroadcastListItem, 0, len(list))
	for _, b := range list {
		result = append(result, BroadcastListItem{
			ID:        b.ID,
			Judul:     b.Judul,
			Deskripsi: b.Deskripsi,
			Tipe:      b.Tipe,
			CabangID:  b.CabangID,
			CreatedAt: b.CreatedAt,
		})
	}
	return result
}

// Delete menghapus satu broadcast berdasarkan ID
func (s *broadcastService) Delete(id uint) error {
	return s.repo.Delete(id)
}

// DeleteAll menghapus seluruh broadcast
func (s *broadcastService) DeleteAll() error {
	return s.repo.DeleteAll()
}
