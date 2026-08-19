package redis

import (
	"context"
	"log"
	"time"
)

type JobHandler func(payload []byte) error

// StartWorker menjalankan listener antrian Redis (Blocking POP)
func StartWorker(ctx context.Context, queueName string, handler JobHandler) {
	log.Printf("👷 Worker mendengarkan antrian: %s", queueName)

	for {
		select {
		case <-ctx.Done():
			log.Println("🛑 Worker dihentikan")
			return
		default:
			// BRPop: Tunggu job masuk max 5 detik per iterasi
			result, err := Client.BRPop(ctx, 5*time.Second, queueName).Result()
			if err != nil {
				// Nil error timeout biasa terjadi jika antrian kosong, lanjut loop
				continue
			}

			if len(result) >= 2 {
				payload := []byte(result[1])
				if err := handler(payload); err != nil {
					log.Printf("❌ Error memproses job: %v", err)
				} else {
					log.Printf("✅ Job berhasil diproses dari antrian %s", queueName)
				}
			}
		}
	}
}
