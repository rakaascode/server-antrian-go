package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rakaascode/server-antrian-go.git/internal/antrian"
	"github.com/rakaascode/server-antrian-go.git/internal/auth"
	"github.com/rakaascode/server-antrian-go.git/internal/broadcast"
	"github.com/rakaascode/server-antrian-go.git/internal/cabang"
	"github.com/rakaascode/server-antrian-go.git/internal/crm"
	"github.com/rakaascode/server-antrian-go.git/internal/user"
	"github.com/rakaascode/server-antrian-go.git/pkg/database"
	"github.com/rakaascode/server-antrian-go.git/routes"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env tidak ditemukan, menggunakan system env")
	}

	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal("❌ Gagal konek DB:", err)
	}
	log.Println("✅ DB connected")

	// Auto migrate
	if err := db.AutoMigrate(
		&cabang.Cabang{},
		&user.User{},
		&antrian.Antrian{},
		&broadcast.Broadcast{},
	); err != nil {
		log.Fatal("❌ Gagal migrasi:", err)
	}
	log.Println("✅ Migrasi selesai")

	// Seed data awal cabang (hanya berjalan jika tabel masih kosong)
	cabang.Seed(db)

	// Wire dependencies
	userRepo := user.NewUserRepository(db)

	authSvc := auth.NewAuthService(userRepo)
	authHandler := auth.NewAuthHandler(authSvc)

	userSvc := user.NewUserService(userRepo)
	userHandler := user.NewUserHandler(userSvc)

	antrianRepo := antrian.NewAntrianRepository(db)
	antrianSvc := antrian.NewAntrianService(antrianRepo)
	antrianHandler := antrian.NewAntrianHandler(antrianSvc)

	crmSvc := crm.NewCrmService(antrianRepo)
	crmHandler := crm.NewCrmHandler(crmSvc)

	broadcastRepo := broadcast.NewBroadcastRepository(db)
	broadcastSvc := broadcast.NewBroadcastService(broadcastRepo, antrianRepo)
	broadcastHandler := broadcast.NewBroadcastHandler(broadcastSvc)

	cabangRepo := cabang.NewCabangRepository(db)
	cabangSvc := cabang.NewCabangService(cabangRepo)
	cabangHandler := cabang.NewCabangHandler(cabangSvc)

	// Setup router
	r := gin.Default()
	routes.SetupRoutes(r, routes.Handlers{
		Auth:      authHandler,
		User:      userHandler,
		Antrian:   antrianHandler,
		CRM:       crmHandler,
		Cabang:    cabangHandler,
		Broadcast: broadcastHandler,
	})

	log.Println("🚀 Server berjalan di :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("❌ Gagal menjalankan server:", err)
	}
}
