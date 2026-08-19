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
	// ── Public Pages ──────────────────────────────────────────────────────────
	r.GET("/privacy-policy", func(c *gin.Context) {
		c.String(200, `Last updated: 14 May 2026

Teduh Service App respects your privacy and is committed to protecting any personal data that may be collected while using our application.

1. Information We Collect
Our application may collect the following information:
- Camera access (if user uses scan or capture features)
- Location data (coarse and precise location) to provide nearby service features
- Network information (WiFi and internet status) for app functionality

We do NOT collect sensitive personal data such as passwords, financial information, or identity documents.

2. How We Use Information
We use the collected information to:
- Provide core app features (such as service queue and location-based services)
- Improve app performance and user experience
- Ensure app functionality over network connections

3. Camera Permission
The camera is only used when the user explicitly activates features that require image capture or scanning. We do not store or upload camera data without user consent.

4. Location Permission
Location data is used only to show relevant nearby services. Location is not shared with third parties.

5. Data Sharing
We do not sell, trade, or share user data with third parties.

6. Data Security
We take reasonable measures to protect user data from unauthorized access or misuse.

7. Children's Privacy
This application is not intended for children under 13. We do not knowingly collect data from children.

8. Changes to This Policy
We may update this Privacy Policy from time to time. Updates will be posted in this page.

9. Contact Us
If you have questions about this Privacy Policy, contact us at:
lteduh-antrean@googlegroups.com`)
	})

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

	// ── Ringkasan semua cabang — harus sebelum /cabang/:id/antrian ────────────
	api.GET("/cabang/antrian/ringkasan", h.Antrian.GetRingkasanSemuaCabang) // 🔵 ringkasan semua cabang hari ini

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
		// Batalkan antrian milik sendiri (hanya bisa jika status "menunggu")
		userProtected.DELETE("/antrian/:id/batal", h.Antrian.BatalkanAntrian)
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
		adminOnly.GET("/crm/antrian", h.CRM.GetAntrianForCrm)        // Picker: list antrian hari ini + antrian_id
		adminOnly.POST("/crm/send", h.CRM.ManualSend)                // Mode 1: manual (input no WA + pesan)
		adminOnly.POST("/crm/reminders", h.CRM.ReminderFromAntrian)  // Mode 2: dari data antrian

		// Broadcast (notifikasi & promo)
		adminOnly.POST("/broadcast", h.Broadcast.Create)             // kirim broadcast baru
		adminOnly.GET("/broadcast/all", h.Broadcast.GetAll)          // lihat semua broadcast yang pernah dikirim
		adminOnly.DELETE("/broadcast/all", h.Broadcast.DeleteAll)    // hapus SEMUA broadcast sekaligus
		adminOnly.DELETE("/broadcast/:id", h.Broadcast.Delete)       // hapus satu broadcast by ID

		// Manajemen Cabang
		adminOnly.POST("/cabang", h.Cabang.Create)
		adminOnly.PUT("/cabang/:id", h.Cabang.Update)
		adminOnly.DELETE("/cabang/:id", h.Cabang.Delete)

		// Manajemen User
		adminOnly.GET("/admin/users/kontak", h.User.GetAllUsersKontak)
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