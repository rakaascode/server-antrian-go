package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/rakaascode/server-antrian-go.git/internal/cabang"
	"github.com/rakaascode/server-antrian-go.git/internal/user"
	"github.com/rakaascode/server-antrian-go.git/pkg/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 1. Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Peringatan: file .env tidak ditemukan, menggunakan variabel environment sistem.")
	}

	// 2. Connect ke Database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("============================================")
	fmt.Println("   SEEDER: Buat Akun Admin Cabang")
	fmt.Println("============================================")

	// 3. Tampilkan daftar cabang
	var cabangs []cabang.Cabang
	if err := db.Order("id asc").Find(&cabangs).Error; err != nil {
		log.Fatalf("Gagal mengambil data cabang: %v", err)
	}

	if len(cabangs) == 0 {
		fmt.Println("\nBelum ada data cabang di database!")
		fmt.Println("Silakan tambahkan cabang terlebih dahulu via API: POST /api/cabang")
		os.Exit(1)
	}

	fmt.Println("\nDaftar Cabang yang tersedia:")
	fmt.Println("--------------------------------------------")
	for _, c := range cabangs {
		fmt.Printf("  [%d] %s — %s\n", c.ID, c.Nama, c.Kota)
	}
	fmt.Println("--------------------------------------------")

	// 4. Pilih cabang
	fmt.Print("\nMasukkan ID Cabang untuk admin ini: ")
	cabangIDStr, _ := reader.ReadString('\n')
	cabangIDStr = strings.TrimSpace(cabangIDStr)

	cabangIDInt, err := strconv.ParseUint(cabangIDStr, 10, 32)
	if err != nil || cabangIDInt == 0 {
		log.Fatalf("ID Cabang tidak valid.")
	}
	cabangIDVal := uint(cabangIDInt)

	// Validasi: pastikan cabang dengan ID tersebut ada
	var selectedCabang cabang.Cabang
	if err := db.First(&selectedCabang, cabangIDVal).Error; err != nil {
		log.Fatalf("Cabang dengan ID %d tidak ditemukan.", cabangIDVal)
	}

	fmt.Printf("\nCabang dipilih: %s (%s)\n", selectedCabang.Nama, selectedCabang.Kota)

	// 5. Input data admin
	fmt.Println("\nIsi data akun admin baru:")
	fmt.Print("Nama Lengkap: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print("Username (untuk login): ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Print("Password (min. 6 karakter): ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	if len(password) < 6 {
		log.Fatalf("Password minimal 6 karakter.")
	}

	// 6. Cek apakah username sudah ada
	var existing user.User
	if err := db.Where("username = ?", username).First(&existing).Error; err == nil {
		log.Fatalf("Username '%s' sudah digunakan. Pilih username lain.", username)
	}

	// 7. Hash Password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		log.Fatalf("Gagal melakukan hashing password: %v", err)
	}

	// 8. Simpan ke Database
	newUser := user.User{
		Name:     name,
		Username: username,
		Password: hashedPassword,
		Role:     "admin",
		CabangID: &cabangIDVal,
	}

	if err := db.Create(&newUser).Error; err != nil {
		log.Fatalf("Gagal membuat user admin: %v", err)
	}

	fmt.Println("\n============================================")
	fmt.Printf("  Admin berhasil dibuat!\n")
	fmt.Printf("  Nama     : %s\n", name)
	fmt.Printf("  Username : %s\n", username)
	fmt.Printf("  Cabang   : %s (%s)\n", selectedCabang.Nama, selectedCabang.Kota)
	fmt.Println("  Role     : admin")
	fmt.Println("============================================")
	fmt.Println("\nGunakan endpoint berikut untuk login:")
	fmt.Println("  POST /api/auth/admin/login")
	fmt.Println("  Body: { \"username\": \"...\", \"password\": \"...\" }")
}
