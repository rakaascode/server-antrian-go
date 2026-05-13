package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rakaascode/server-antrian-go.git/internal/antrian"
	"github.com/rakaascode/server-antrian-go.git/internal/auth"
	"github.com/rakaascode/server-antrian-go.git/internal/broadcast"
	"github.com/rakaascode/server-antrian-go.git/internal/cabang"
	"github.com/rakaascode/server-antrian-go.git/internal/crm"
	"github.com/rakaascode/server-antrian-go.git/internal/user"
	"github.com/rakaascode/server-antrian-go.git/pkg/middleware"
)

type Handlers struct {
	Auth      *auth.AuthHandler
	User      *user.UserHandler
	Antrian   *antrian.AntrianHandler
	CRM       *crm.CrmHandler
	Cabang    *cabang.CabangHandler
	Broadcast *broadcast.BroadcastHandler
}

func SetupRoutes(r *gin.Engine, h Handlers) {
	api := r.Group("/api")

	// ── Auth ─────────────────────────────────────────────────────────────────
	authGroup := api.Group("/auth")
	{
		authGroup.POST("/google", h.Auth.GoogleLogin)       // user Android (Google Sign-In)
		authGroup.POST("/admin/login", h.Auth.AdminLogin)   // admin cabang (username + password)
	}

	// ── Cabang (public) ───────────────────────────────────────────────────────
	api.GET("/cabang", h.Cabang.GetAll)
	api.GET("/cabang/:id", h.Cabang.GetByID)

	// ── Antrian per cabang (public — data non-sensitif) ───────────────────────
	api.GET("/cabang/:id/antrian", h.Antrian.GetByCabang)
	api.GET("/cabang/:id/antrian/status", h.Antrian.GetStatusCabang) // 🔴 nomor yang sedang dipanggil


	// ── Protected: User ───────────────────────────────────────────────────────
	userProtected := api.Group("/")
	userProtected.Use(middleware.AuthMiddleware())
	{
		// Profil user sendiri (nama, email, avatar, riwayat antrian + cabang)
		userProtected.GET("/user/profile", h.User.GetProfile)
		userProtected.GET("/users/profile", h.User.GetProfile) // alias plural

		// Update profil (alamat, avatar, nama)
		userProtected.PUT("/user/profile", h.User.UpdateProfile)
		userProtected.PUT("/users/profile", h.User.UpdateProfile) // alias plural

		// Kontak WA user — simpan, lihat, update, hapus
		userProtected.GET("/users/kontak", h.User.GetKontak)
		userProtected.POST("/users/kontak", h.User.SaveKontak)
		userProtected.PUT("/users/kontak", h.User.SaveKontak) // alias update
		userProtected.DELETE("/users/kontak", h.User.DeleteKontak)

		// Ambil nomor antrian (wajib login)
		userProtected.POST("/antrian", h.Antrian.AmbilAntrian)
		// Detail antrian milik sendiri (owner atau admin)
		userProtected.GET("/antrian/:id", h.Antrian.GetByID)
		// Posisi di antrian (berapa orang di depan)
		userProtected.GET("/antrian/:id/posisi", h.Antrian.GetPosisi)
		// Semua antrian milik user ini
		userProtected.GET("/antrian/me", h.Antrian.GetMyAntrian)

		// Broadcast: list (promo + antrian cabang user) + detail
		userProtected.GET("/broadcast", h.Broadcast.GetForUser)
		userProtected.GET("/broadcast/:id", h.Broadcast.GetByID)
	}

	// ── Protected: Admin cabang ───────────────────────────────────────────────
	adminOnly := api.Group("/")
	adminOnly.Use(middleware.AuthMiddleware(), middleware.RequireRole("admin"))
	{
		// Kelola cabang sendiri
		adminOnly.GET("/cabang/:id/antrian/detail", h.Antrian.GetByCabangAdmin)
		adminOnly.POST("/antrian/call-next", h.Antrian.CallNext)
		adminOnly.PUT("/antrian/:id/selesai", h.Antrian.Selesai)
		adminOnly.DELETE("/antrian/:id", h.Antrian.Delete)

		// CRM / Pengingat WA
		adminOnly.POST("/crm/send", h.CRM.ManualSend)               // Mode 1: manual (input no WA + pesan)
		adminOnly.POST("/crm/reminders", h.CRM.ReminderFromAntrian)  // Mode 2: dari data antrian

		// Broadcast (notifikasi & promo)
		adminOnly.POST("/broadcast", h.Broadcast.Create)          // kirim broadcast baru
		adminOnly.GET("/broadcast/all", h.Broadcast.GetAll)       // lihat semua broadcast yang pernah dikirim

		// Manajemen Cabang
		adminOnly.POST("/cabang", h.Cabang.Create)
		adminOnly.PUT("/cabang/:id", h.Cabang.Update)
		adminOnly.DELETE("/cabang/:id", h.Cabang.Delete)

		// Manajemen User
		adminOnly.GET("/users", h.User.GetAll)
		adminOnly.GET("/users/:id", h.User.GetByID)
		adminOnly.POST("/users", h.User.Create)
		adminOnly.PUT("/users/:id", h.User.Update)
		adminOnly.DELETE("/users/:id", h.User.Delete)
	}

	// ── Protected: Super Admin (global admin tanpa cabang_id) ─────────────────
	superAdmin := api.Group("/super")
	superAdmin.Use(middleware.AuthMiddleware(), middleware.RequireSuperAdmin())
	{
		// Buat akun admin cabang baru
		superAdmin.POST("/admin/cabang", h.User.CreateAdminCabang)

		// List semua admin (seluruh cabang)
		superAdmin.GET("/admins", h.User.GetAllAdmins)

		// List admin di cabang tertentu
		superAdmin.GET("/cabang/:id/admins", h.User.GetAdminsByCabang)

		// Tugaskan admin ke cabang (atau ganti cabang)
		superAdmin.PUT("/admins/:id/assign", h.User.AssignCabang)

		// Lepas admin dari cabang (cabang_id = NULL)
		superAdmin.DELETE("/admins/:id/assign", h.User.UnassignCabang)
	}
}