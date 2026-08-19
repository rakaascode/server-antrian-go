package redis

import (
	"context"
	"fmt"
	"log"
)

const NotificationQueue = "queue:notifications"

// PushNotification memasukkan job notifikasi ke antrian Redis
func PushNotification(ctx context.Context, payload []byte) error {
	if Client == nil {
		return fmt.Errorf("redis client belum diinisialisasi")
	}

	err := Client.LPush(ctx, NotificationQueue, payload).Err()
	if err != nil {
		log.Printf("❌ Gagal push job ke Redis: %v", err)
		return err
	}
	return nil
}
