package notification

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"gorm.io/gorm"
)

// ScheduledAntrianItem merepresentasikan data antrian untuk pengingat otomatis
type ScheduledAntrianItem struct {
	ID                uint
	CabangID          uint
	NamaPemilik       string
	NoWAReminder      string
	MerkMotor         string
	TipeMotor         string
	NomorAntrian      int
	TanggalKedatangan time.Time
	EstimasiJam       string
}

// StartAutoReminderScheduler menjalankan background cron scheduler tiap 1 menit
// untuk mendeteksi antrian hari ini yang berjarak ~30 menit dari jadwal kedatangan,
// lalu mengirimkan WhatsApp reminder otomatis via Fonnte.
func StartAutoReminderScheduler(ctx context.Context, db *gorm.DB) {
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		defer ticker.Stop()
		log.Println("⏰ Auto WA Reminder Scheduler (H-30 Menit) Aktif!")

		for {
			select {
			case <-ctx.Done():
				log.Println("🛑 Auto WA Reminder Scheduler dihentikan.")
				return
			case <-ticker.C:
				checkAndSendUpcomingReminders(db)
			}
		}
	}()
}

func checkAndSendUpcomingReminders(db *gorm.DB) {
	now := time.Now()
	targetTime := now.Add(30 * time.Minute)
	targetHourMinute := targetTime.Format("15:04") // e.g. "09:30"
	targetHour := targetTime.Format("15:00")       // e.g. "09:00"

	var antrians []ScheduledAntrianItem
	err := db.Table("antrians").
		Select("id, cabang_id, nama_pemilik, no_wa_reminder, merk_motor, tipe_motor, nomor_antrian, tanggal_kedatangan, estimasi_jam").
		Where("status = ? AND reminder_aktif = ? AND no_wa_reminder != ? AND DATE(tanggal_kedatangan) = CURRENT_DATE", "menunggu", true, "").
		Scan(&antrians).Error

	if err != nil {
		log.Printf("⚠️ Scheduler Error fetching antrian: %v\n", err)
		return
	}

	for _, a := range antrians {
		cleanEst := strings.TrimSpace(a.EstimasiJam)
		// Cocokkan jika estimasi jam dimulai dari jam target (misal: "09:00 - 10:00" atau "09:30")
		if strings.HasPrefix(cleanEst, targetHourMinute) || strings.HasPrefix(cleanEst, targetHour) {
			nomorDisplay := fmt.Sprintf("A-%03d", a.NomorAntrian)
			pesan := fmt.Sprintf(
				"⏰ *Pengingat Servis H-30 Menit*\n\nHalo %s!\n\nKami mengingatkan jadwal servis motor Anda akan dimulai dalam 30 menit ke depan:\n\n🏍️ Motor: %s %s\n🔢 Nomor Antrean: *#%s*\n⏰ Estimasi Jam: %s\n\nMohon bersiap dan datang tepat waktu ke cabang bengkel. Terima kasih! 🙏",
				a.NamaPemilik,
				a.MerkMotor,
				a.TipeMotor,
				nomorDisplay,
				a.EstimasiJam,
			)

			// Kirim WA secara async
			go func(target, msg string, id uint) {
				if err := SendWA(target, msg); err != nil {
					log.Printf("❌ Gagal kirim WA H-30 ke antrian ID %d: %v\n", id, err)
				} else {
					log.Printf("✅ Berhasil kirim WA H-30 ke antrian ID %d (No: %s)\n", id, target)
				}
			}(a.NoWAReminder, pesan, a.ID)
		}
	}
}
