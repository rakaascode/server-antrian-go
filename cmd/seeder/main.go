package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/rakaascode/server-antrian-go.git/internal/user"
	"github.com/rakaascode/server-antrian-go.git/pkg/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Daftar admin yang akan di-seed.
// Ubah username/password sesuai kebutuhan sebelum menjalankan.
var adminList = []struct {
	Name     string
	Username string
	Password string
}{
	{Name: "Admin Antrian", Username: "admin_antrian", Password: "Admin@Antrian123"},
	{Name: "Admin CRM", Username: "admin_crm", Password: "Admin@CRM123"},
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Peringatan: .env tidak ditemukan, menggunakan environment variabel.")
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"), os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal koneksi ke database: %v", err)
	}

	fmt.Println("============================================")
	fmt.Println("   SEEDER: Membuat Akun Admin")
	fmt.Println("============================================")

	for _, a := range adminList {
		// Skip jika username sudah ada
		var existing user.User
		if err := db.Where("username = ?", a.Username).First(&existing).Error; err == nil {
			fmt.Printf("  [SKIP] Username '%s' sudah ada di database.\n", a.Username)
			continue
		}

		hashed, err := utils.HashPassword(a.Password)
		if err != nil {
			log.Printf("  [ERROR] Gagal hash password untuk %s: %v\n", a.Username, err)
			continue
		}

		newUser := user.User{
			Name:     a.Name,
			Username: &a.Username,
			Password: hashed,
			Role:     "admin",
		}

		if err := db.Create(&newUser).Error; err != nil {
			log.Printf("  [ERROR] Gagal membuat user %s: %v\n", a.Username, err)
			continue
		}

		fmt.Printf("  [OK] %-20s | username: %-20s | password: %s\n", a.Name, a.Username, a.Password)
	}

	fmt.Println("============================================")
	fmt.Println("Selesai! Login via: POST /api/auth/admin/login")
}
