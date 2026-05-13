package antrian

import (
	"errors"
	"fmt"

	"github.com/rakaascode/server-antrian-go.git/internal/notification"
)

type AntrianService interface {
	GetByCabang(cabangID uint, status string) ([]Antrian, error)
	GetByID(id uint) (Antrian, error)
	GetMyAntrian(userID uint) ([]Antrian, error)
	AmbilAntrian(req AmbilAntrianRequest, userID *uint) (Antrian, error)
	CallNext(cabangID uint) (Antrian, error)
	Selesai(id, cabangID uint) error
	Delete(id, cabangID uint) error
	GetStatusCabang(cabangID uint) (StatusCabangResponse, error)
	GetPosisi(antrianID uint) (PosisiResponse, error)
}

type antrianService struct {
	repo AntrianRepository
}

func NewAntrianService(r AntrianRepository) AntrianService {
	return &antrianService{r}
}

func (s *antrianService) GetByCabang(cabangID uint, status string) ([]Antrian, error) {
	if status != "" {
		return s.repo.FindByCabangAndStatus(cabangID, status)
	}
	return s.repo.FindByCabang(cabangID)
}

func (s *antrianService) GetByID(id uint) (Antrian, error) {
	return s.repo.FindByID(id)
}

func (s *antrianService) GetMyAntrian(userID uint) ([]Antrian, error) {
	return s.repo.FindByUserID(userID)
}

// AmbilAntrian membuat nomor antrian baru secara otomatis per cabang
func (s *antrianService) AmbilAntrian(req AmbilAntrianRequest, userID *uint) (Antrian, error) {
	// Validasi: jika reminder aktif, no WA wajib diisi
	if req.ReminderAktif && req.NoWAReminder == "" {
		return Antrian{}, errors.New("no_wa_reminder wajib diisi jika reminder_aktif = true")
	}

	count, err := s.repo.CountTodayByCabang(req.CabangID)
	if err != nil {
		return Antrian{}, err
	}

	a := Antrian{
		CabangID:          req.CabangID,
		UserID:            userID,
		NomorAntrian:      int(count) + 1,
		Status:            StatusMenunggu,
		NamaPemilik:       req.NamaPemilik,
		NoHP:              req.NoHP,
		MerkMotor:         req.MerkMotor,
		TipeMotor:         req.TipeMotor,
		NoRangka:          req.NoRangka,
		NoMesin:           req.NoMesin,
		TahunPembuatan:    req.TahunPembuatan,
		TanggalKedatangan: req.TanggalKedatangan,
		EstimasiJam:       req.EstimasiJam,
		ReminderAktif:     req.ReminderAktif,
		NoWAReminder:      req.NoWAReminder,
		Catatan:           req.Catatan,
	}
	return s.repo.Create(a)
}

// CallNext panggil antrian berikutnya di cabang ini
// Jika antrian memiliki reminder aktif, kirim notifikasi WA secara async
func (s *antrianService) CallNext(cabangID uint) (Antrian, error) {
	list, err := s.repo.FindByCabangAndStatus(cabangID, StatusMenunggu)
	if err != nil || len(list) == 0 {
		return Antrian{}, errors.New("tidak ada antrian yang menunggu di cabang ini")
	}
	next := list[0]
	if err := s.repo.UpdateStatus(next.ID, StatusDipanggil); err != nil {
		return Antrian{}, err
	}
	next.Status = StatusDipanggil

	// Auto-kirim WA reminder jika user mengaktifkan pengingat
	if next.ReminderAktif && next.NoWAReminder != "" {
		pesan := fmt.Sprintf(
			"🔔 *Pengingat Antrian*\n\nHalo %s! Nomor antrian Anda *#%d* sedang dipanggil. Silakan segera menuju ke area servis.\n\n🏍️ %s %s\nEstimasi: %s\n\nTerima kasih! 🙏",
			next.NamaPemilik,
			next.NomorAntrian,
			next.MerkMotor,
			next.TipeMotor,
			next.EstimasiJam,
		)
		// Kirim secara async agar tidak menghambat response
		go notification.SendWA(next.NoWAReminder, pesan)
	}

	return next, nil
}

// Selesai tandai antrian selesai — validasi cabang
func (s *antrianService) Selesai(id, cabangID uint) error {
	a, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("antrian tidak ditemukan")
	}
	if cabangID != 0 && a.CabangID != cabangID {
		return errors.New("antrian bukan milik cabang Anda")
	}

	if err := s.repo.UpdateStatus(id, StatusSelesai); err != nil {
		return err
	}

	// Notifikasi selesai jika reminder aktif
	if a.ReminderAktif && a.NoWAReminder != "" {
		pesan := fmt.Sprintf(
			"✅ *Servis Selesai*\n\nHalo %s! Motor %s %s Anda telah selesai diservis. Silakan menuju kasir. Terima kasih telah mempercayai kami! 🙏",
			a.NamaPemilik,
			a.MerkMotor,
			a.TipeMotor,
		)
		go notification.SendWA(a.NoWAReminder, pesan)
	}

	return nil
}

// Delete hapus antrian — validasi cabang
func (s *antrianService) Delete(id, cabangID uint) error {
	a, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("antrian tidak ditemukan")
	}
	if cabangID != 0 && a.CabangID != cabangID {
		return errors.New("antrian bukan milik cabang Anda")
	}
	return s.repo.Delete(id)
}

// GetStatusCabang info realtime: nomor yang sedang dipanggil + total menunggu
func (s *antrianService) GetStatusCabang(cabangID uint) (StatusCabangResponse, error) {
	totalMenunggu, err := s.repo.CountTodayByCabang(cabangID)
	if err != nil {
		return StatusCabangResponse{}, err
	}

	// Hitung hanya yang berstatus "menunggu" hari ini
	listMenunggu, _ := s.repo.FindByCabangAndStatus(cabangID, StatusMenunggu)
	total := int64(len(listMenunggu))

	dipanggil, err := s.repo.FindLatestDipanggil(cabangID)
	if err != nil {
		// Belum ada yang dipanggil hari ini
		_ = totalMenunggu
		return StatusCabangResponse{
			NomorDipanggil: nil,
			StatusPanggil:  "belum ada",
			TotalMenunggu:  total,
		}, nil
	}

	nomor := dipanggil.NomorAntrian
	return StatusCabangResponse{
		NomorDipanggil: &nomor,
		StatusPanggil:  StatusDipanggil,
		TotalMenunggu:  total,
	}, nil
}

// GetPosisi menghitung posisi user di antrian: berapa orang di depannya
func (s *antrianService) GetPosisi(antrianID uint) (PosisiResponse, error) {
	a, err := s.repo.FindByID(antrianID)
	if err != nil {
		return PosisiResponse{}, errors.New("antrian tidak ditemukan")
	}

	// Jika sudah dipanggil atau selesai
	if a.Status != StatusMenunggu {
		var pesan string
		if a.Status == StatusDipanggil {
			pesan = "Nomor antrian Anda sedang dipanggil! Segera ke loket."
		} else {
			pesan = "Servis Anda sudah selesai."
		}
		return PosisiResponse{
			NomorAntrian: a.NomorAntrian,
			Status:       a.Status,
			Posisi:       0,
			Pesan:        pesan,
		}, nil
	}

	// Hitung antrian "menunggu" dengan nomor lebih kecil (= orang di depan)
	posisi, err := s.repo.CountMenungguSebelum(a.CabangID, a.NomorAntrian)
	if err != nil {
		return PosisiResponse{}, err
	}

	// Nomor yang sedang dilayani
	dipanggil, _ := s.repo.FindLatestDipanggil(a.CabangID)
	var nomorDipanggil *int
	if dipanggil != nil {
		n := dipanggil.NomorAntrian
		nomorDipanggil = &n
	}

	var pesan string
	if posisi == 0 {
		pesan = "Anda adalah antrian berikutnya!"
	} else if posisi == 1 {
		pesan = "1 orang lagi di depan Anda."
	} else {
		pesan = fmt.Sprintf("%d orang di depan Anda.", posisi)
	}

	return PosisiResponse{
		NomorAntrian:   a.NomorAntrian,
		Status:         a.Status,
		Posisi:         posisi,
		NomorDipanggil: nomorDipanggil,
		Pesan:          pesan,
	}, nil
}
