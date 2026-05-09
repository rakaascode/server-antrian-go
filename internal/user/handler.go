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
