package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service UserService
}

func NewUserHandler(s UserService) *UserHandler {
	return &UserHandler{s}
}

// GetProfile GET /api/user/profile  (wajib login — user Android)
// Mengembalikan data profil user dari JWT + riwayat antrian dengan info cabang
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	profile, err := h.service.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "User tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": profile})
}

// UpdateProfile PUT /api/user/profile  (wajib login)
// Memperbarui profil user termasuk avatar dan alamat
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	u, err := h.service.UpdateProfile(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	u.Password = ""
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Profil berhasil diperbarui",
		"data":    u,
	})
}

func (h *UserHandler) GetAll(c *gin.Context) {
	users, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	user, err := h.service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) Create(c *gin.Context) {
	var u User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.Create(u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *UserHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var u User
	c.ShouldBindJSON(&u)

	result, err := h.service.Update(uint(id), u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *UserHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	err := h.service.Delete(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// GetKontak GET /api/users/kontak — lihat nomor WA yang tersimpan
func (h *UserHandler) GetKontak(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	u, err := h.service.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "User tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"user_id": u.ID,
			"no_wa":   u.NoWA,
		},
	})
}

// SaveKontak POST /api/users/kontak — simpan atau update nomor WA
func (h *UserHandler) SaveKontak(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var req KontakRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "no_wa wajib diisi"})
		return
	}

	u, err := h.service.SaveKontak(userID, req.NoWA)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Kontak WA berhasil disimpan",
		"data": gin.H{
			"user_id": u.ID,
			"no_wa":   u.NoWA,
		},
	})
}

// DeleteKontak DELETE /api/users/kontak — hapus nomor WA
func (h *UserHandler) DeleteKontak(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	u, err := h.service.DeleteKontak(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Kontak WA berhasil dihapus",
		"data": gin.H{
			"user_id": u.ID,
			"no_wa":   "",
		},
	})
}

// ──────────────────────────────────────────────────────────────────────────────
// Admin Cabang — hanya bisa diakses super admin
// ──────────────────────────────────────────────────────────────────────────────

// CreateAdminCabang POST /api/super/admin/cabang — buat akun admin cabang baru
func (h *UserHandler) CreateAdminCabang(c *gin.Context) {
	var req CreateAdminCabangRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	u, err := h.service.CreateAdminCabang(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	u.Password = ""
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Admin cabang berhasil dibuat",
		"data":    u,
	})
}

// GetAllAdmins GET /api/super/admins — list semua admin (semua cabang)
func (h *UserHandler) GetAllAdmins(c *gin.Context) {
	admins, err := h.service.GetAllAdmins()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	for i := range admins {
		admins[i].Password = ""
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": admins})
}

// GetAdminsByCabang GET /api/super/cabang/:id/admins — list admin di cabang tertentu
func (h *UserHandler) GetAdminsByCabang(c *gin.Context) {
	cabangID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID cabang tidak valid"})
		return
	}

	admins, err := h.service.GetAdminsByCabang(uint(cabangID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	for i := range admins {
		admins[i].Password = ""
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": admins})
}

// AssignCabang PUT /api/super/admins/:id/assign — tugaskan admin ke cabang
func (h *UserHandler) AssignCabang(c *gin.Context) {
	adminID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID admin tidak valid"})
		return
	}

	var body struct {
		CabangID uint `json:"cabang_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "cabang_id wajib diisi"})
		return
	}

	u, err := h.service.AssignCabang(uint(adminID), body.CabangID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	u.Password = ""
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Admin berhasil ditugaskan ke cabang",
		"data":    u,
	})
}

// UnassignCabang DELETE /api/super/admins/:id/assign — lepas admin dari cabang
func (h *UserHandler) UnassignCabang(c *gin.Context) {
	adminID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID admin tidak valid"})
		return
	}

	u, err := h.service.UnassignCabang(uint(adminID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	u.Password = ""
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Admin berhasil dilepas dari cabang",
		"data":    u,
	})
}

// GetAllUsersKontak GET /api/admin/users/kontak — admin mengambil semua kontak user
func (h *UserHandler) GetAllUsersKontak(c *gin.Context) {
	users, err := h.service.GetAllUsersKontak()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	type KontakResponse struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
		NoWA string `json:"no_wa"`
	}

	var data []KontakResponse
	for _, u := range users {
		data = append(data, KontakResponse{
			ID:   u.ID,
			Name: u.Name,
			NoWA: u.NoWA,
		})
	}

	// Pastikan return data slice kosong jika tidak ada user, bukan null
	if data == nil {
		data = []KontakResponse{}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}
