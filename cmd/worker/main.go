package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/rakaascode/server-antrian-go.git/internal/notification"
	"github.com/rakaascode/server-antrian-go.git/pkg/redis"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env tidak ditemukan, menggunakan system env")
	}

	// Inisialisasi koneksi Redis
	_, err := redis.InitRedis()
	if err != nil {
		log.Fatalf("❌ Gagal terhubung ke Redis: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Menangkap signal shutdown (Graceful Shutdown)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("⚠️ Menerima sinyal stop, mematikan worker...")
		cancel()
	}()

	log.Println("🚀 Background Worker Notifikasi Aktif!")

	// Jalankan consumer antrian notifikasi
	redis.StartWorker(ctx, redis.NotificationQueue, func(payload []byte) error {
		job, err := notification.ParseJob(payload)
		if err != nil {
			return err
		}
		return job.Process()
	})
}
