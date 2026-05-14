package antrian

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AntrianHandler struct {
	service AntrianService
}

func NewAntrianHandler(s AntrianService) *AntrianHandler {
	return &AntrianHandler{s}
}

// GetByCabang GET /api/cabang/:id/antrian?status=menunggu
// Public: hanya tampilkan data non-sensitif
func (h *AntrianHandler) GetByCabang(c *gin.Context) {
	cabangID, _ := strconv.Atoi(c.Param("id"))
	status := c.Query("status")

	list, err := h.service.GetByCabang(uint(cabangID), status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	// Tampilkan versi publik (tanpa data sensitif)
	var pub []AntrianPublicResponse
	for _, a := range list {
		pub = append(pub, AntrianPublicResponse{
			ID:                a.ID,
			NomorAntrian:      a.NomorAntrian,
			Status:            a.Status,
			EstimasiJam:       a.EstimasiJam,
			TanggalKedatangan: a.TanggalKedatangan,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"success":       true,
		"total_antrian": len(list),
		"data":          pub,
	})
}

// GetByCabangAdmin GET /api/cabang/:id/antrian/detail  (admin only)
// Menampilkan data lengkap antrian di cabangnya (atau semua cabang jika Super Admin)
func (h *AntrianHandler) GetByCabangAdmin(c *gin.Context) {
	cabangIDURL, _ := strconv.Atoi(c.Param("id"))

	cabangIDVal, exists := c.Get("cabang_id")
	if exists {
		// Admin cabang biasa
		if uint(cabangIDURL) != cabangIDVal.(uint) {
			c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Akses ditolak: Admin tidak terikat ke cabang ini"})
			return
		}
	}
	// Jika tidak exists (tapi lolos RequireRole admin di route), berarti Super Admin

	status := c.Query("status")
	list, err := h.service.GetByCabang(uint(cabangIDURL), status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "total": len(list), "data": list})
}

// GetByID GET /api/antrian/:id  (owner atau admin cabang)
func (h *AntrianHandler) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	a, err := h.service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Antrian tidak ditemukan"})
		return
	}

	// Cek akses: harus owner atau admin cabang yang sama
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")
	cabangID, _ := c.Get("cabang_id")

	isOwner := a.UserID != nil && *a.UserID == userID.(uint)
	isAdminSameCabang := role == "admin" && cabangID != nil && a.CabangID == cabangID.(uint)
	isSuperAdmin := role == "admin" && cabangID == nil

	if !isOwner && !isAdminSameCabang && !isSuperAdmin {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Akses ditolak"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": a})
}

// GetMyAntrian GET /api/antrian/me  (user login)
func (h *AntrianHandler) GetMyAntrian(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	list, err := h.service.GetMyAntrian(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": list})
}

// AmbilAntrian POST /api/antrian  (wajib login)
func (h *AntrianHandler) AmbilAntrian(c *gin.Context) {
	var req AmbilAntrianRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	// UserID selalu ada karena endpoint ini wajib login
	uid := c.MustGet("user_id").(uint)
	userID := &uid

	a, err := h.service.AmbilAntrian(req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Antrian berhasil diambil",
		"data":    a,
	})
}

// CallNext POST /api/antrian/call-next  (admin cabang atau super admin)
func (h *AntrianHandler) CallNext(c *gin.Context) {
	var targetCabangID uint
	cabangIDVal, exists := c.Get("cabang_id")
	if exists {
		targetCabangID = cabangIDVal.(uint)
	} else {
		// Super Admin: harus mengirim cabang_id via query parameter
		cID, _ := strconv.Atoi(c.Query("cabang_id"))
		if cID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Super Admin harus menyertakan ?cabang_id=..."})
			return
		}
		targetCabangID = uint(cID)
	}

	a, err := h.service.CallNext(targetCabangID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Antrian dipanggil", "data": a})
}

// Selesai PUT /api/antrian/:id/selesai  (admin cabang atau super admin)
func (h *AntrianHandler) Selesai(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	
	var adminCabangID uint
	cabangIDVal, exists := c.Get("cabang_id")
	if exists {
		adminCabangID = cabangIDVal.(uint)
	} else {
		adminCabangID = 0 // 0 menandakan Super Admin
	}

	if err := h.service.Selesai(uint(id), adminCabangID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Antrian selesai"})
}

// Delete DELETE /api/antrian/:id  (admin cabang atau super admin)
func (h *AntrianHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	
	var adminCabangID uint
	cabangIDVal, exists := c.Get("cabang_id")
	if exists {
		adminCabangID = cabangIDVal.(uint)
	} else {
		adminCabangID = 0 // 0 menandakan Super Admin
	}

	if err := h.service.Delete(uint(id), adminCabangID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Antrian dihapus"})
}

// GetStatusCabang GET /api/cabang/:id/antrian/status  (public)
// Nomor yang sedang dipanggil + total menunggu di cabang
func (h *AntrianHandler) GetStatusCabang(c *gin.Context) {
	cabangID, _ := strconv.Atoi(c.Param("id"))
	resp, err := h.service.GetStatusCabang(uint(cabangID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

// GetRingkasanSemuaCabang GET /api/cabang/antrian/ringkasan  (public)
// Ringkasan antrian hari ini: nama cabang, nomor dipanggil, estimasi jam, sisa antrian
func (h *AntrianHandler) GetRingkasanSemuaCabang(c *gin.Context) {
	data, err := h.service.GetRingkasanSemuaCabang()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"total_cabang": len(data),
		"data":         data,
	})
}

// GetPosisi GET /api/antrian/:id/posisi  (user login)
// Posisi user di antrian: berapa orang di depan + nomor yang sedang dilayani
func (h *AntrianHandler) GetPosisi(c *gin.Context) {
	antrianID, _ := strconv.Atoi(c.Param("id"))

	a, err := h.service.GetByID(uint(antrianID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Antrian tidak ditemukan"})
		return
	}

	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")
	cabangID, _ := c.Get("cabang_id")

	isOwner := a.UserID != nil && *a.UserID == userID.(uint)
	isAdminSameCabang := role == "admin" && cabangID != nil && a.CabangID == cabangID.(uint)
	isSuperAdmin := role == "admin" && cabangID == nil

	if !isOwner && !isAdminSameCabang && !isSuperAdmin {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Akses ditolak"})
		return
	}

	resp, err := h.service.GetPosisi(uint(antrianID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}
