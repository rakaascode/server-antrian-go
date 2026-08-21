package crm

import (
	"errors"
	"fmt"

	"github.com/rakaascode/server-antrian-go.git/internal/antrian"
	"github.com/rakaascode/server-antrian-go.git/internal/notification"
)

type CrmService interface {
	// Mode 1: Manual — admin input no WA + pesan bebas
	ManualSend(req ManualSendRequest) error

	// Mode 2: Auto dari Antrian — ambil no WA dari data antrian + kirim pesan
	ReminderFromAntrian(req ReminderFromAntrianRequest, cabangID uint) error

	// Listing antrian hari ini di cabang admin (untuk picker sebelum kirim reminder)
	GetAntrianForCrm(cabangID uint) ([]AntrianCrmItem, error)
}

type crmService struct {
	antrianRepo antrian.AntrianRepository
}

func NewCrmService(antrianRepo antrian.AntrianRepository) CrmService {
	return &crmService{antrianRepo: antrianRepo}
}

// ManualSend: admin input no_wa + pesan langsung
func (s *crmService) ManualSend(req ManualSendRequest) error {
	err := notification.SendWA(req.NoWA, req.Pesan)
	if err != nil {
		return fmt.Errorf("gagal mengirim WA ke %s: %v", req.NoWA, err)
	}
	return nil
}

// ReminderFromAntrian: ambil no WA dari data antrian, kirim pesan pengingat
func (s *crmService) ReminderFromAntrian(req ReminderFromAntrianRequest, cabangID uint) error {
	a, err := s.antrianRepo.FindByID(req.AntrianID)
	if err != nil {
		return errors.New("antrian tidak ditemukan")
	}

	// Validasi: admin hanya bisa kirim ke antrian di cabangnya
	if a.CabangID != cabangID {
		return errors.New("antrian bukan dari cabang Anda")
	}

	// Tentukan nomor WA tujuan (prioritas: NoWAReminder, fallback: NoHP)
	noWA := a.NoWAReminder
	if noWA == "" {
		// noWA = a.NoHP (no_polisi bukan nomor wa)
	}
	if noWA == "" {
		return errors.New("nomor WhatsApp pelanggan tidak tersedia di data antrian ini")
	}

	// Gunakan pesan custom atau template default
	pesan := req.Pesan
	if pesan == "" {
		pesan = fmt.Sprintf(
			"📢 *Pengingat dari Bengkel*\n\nHalo %s!\n\nKami mengingatkan Anda memiliki antrian servis motor:\n🏍️ %s %s (%d)\n🔢 Nomor Antrian: *#%d*\n📅 Tanggal: %s\n⏰ Estimasi Jam: %s\n\nMohon hadir tepat waktu. Terima kasih! 🙏",
			a.NamaPemilik,
			a.MerkMotor,
			a.TipeMotor,
			a.TahunPembuatan,
			a.NomorAntrian,
			a.TanggalKedatangan.Format("02 Jan 2006"),
			a.EstimasiJam,
		)
	} else {
		// Inject nama pelanggan di awal pesan custom
		pesan = fmt.Sprintf("Halo %s! %s", a.NamaPemilik, pesan)
	}

	err = notification.SendWA(noWA, pesan)
	if err != nil {
		return fmt.Errorf("gagal mengirim WA: %v", err)
	}
	return nil
}

// GetAntrianForCrm mengembalikan list antrian hari ini yang masih aktif (menunggu/dipanggil)
// di cabang admin — digunakan sebagai picker sebelum admin memilih antrian mana yang mau dikirimi reminder
func (s *crmService) GetAntrianForCrm(cabangID uint) ([]AntrianCrmItem, error) {
	list, err := s.antrianRepo.FindTodayActiveByCabang(cabangID)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data antrian: %v", err)
	}

	result := make([]AntrianCrmItem, 0, len(list))
	for _, a := range list {
		result = append(result, AntrianCrmItem{
			ID:            a.ID,
			NomorAntrian:  a.NomorAntrian,
			Status:        a.Status,
			NamaPemilik:   a.NamaPemilik,
			NoPolisi:      a.NoPolisi,
			NoWAReminder:  a.NoWAReminder,
			ReminderAktif: a.ReminderAktif,
			MerkMotor:     a.MerkMotor,
			TipeMotor:     a.TipeMotor,
			EstimasiJam:   a.EstimasiJam,
			Tanggal:       a.TanggalKedatangan,
		})
	}
	return result, nil
}
