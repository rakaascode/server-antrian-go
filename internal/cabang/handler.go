package cabang

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rakaascode/server-antrian-go.git/internal/antrian"
)

type CabangHandler struct {
	service     CabangService
	antrianSvc  antrian.AntrianService
}

func NewCabangHandler(s CabangService, antrianSvc antrian.AntrianService) *CabangHandler {
	return &CabangHandler{service: s, antrianSvc: antrianSvc}
}

// GetAll GET /api/cabang
// Mengembalikan daftar semua cabang beserta ringkasan antrian hari ini
func (h *CabangHandler) GetAll(c *gin.Context) {
	list, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	// Ambil ringkasan antrian semua cabang, build map untuk lookup O(1)
	ringkasanMap := make(map[uint]antrian.RingkasanCabangResponse)
	if ringkasan, err := h.antrianSvc.GetRingkasanSemuaCabang(); err == nil {
		for _, r := range ringkasan {
			ringkasanMap[r.CabangID] = r
		}
	}

	// Gabungkan data cabang + ringkasan antrian
	result := make([]CabangWithAntrianResponse, 0, len(list))
	for _, cb := range list {
		r := ringkasanMap[cb.ID]
		result = append(result, CabangWithAntrianResponse{
			ID:        cb.ID,
			Nama:      cb.Nama,
			Alamat:    cb.Alamat,
			Kota:      cb.Kota,
			NoTelp:    cb.NoTelp,
			Latitude:  cb.Latitude,
			Longitude: cb.Longitude,
			AntrianHariIni: AntrianHariIni{
				NomorDipanggil: r.NomorDipanggil,
				EstimasiJam:    r.EstimasiJam,
				SisaAntrian:    r.SisaAntrian,
			},
		})
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": result})
}

// GetByID GET /api/cabang/:id
func (h *CabangHandler) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	cb, err := h.service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Cabang tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": cb})
}

// Create POST /api/cabang
func (h *CabangHandler) Create(c *gin.Context) {
	var req CreateCabangRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	cb, err := h.service.Create(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Cabang berhasil dibuat", "data": cb})
}

// Update PUT /api/cabang/:id
func (h *CabangHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req UpdateCabangRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	cb, err := h.service.Update(uint(id), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": cb})
}

// Delete DELETE /api/cabang/:id
func (h *CabangHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.service.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Cabang dihapus"})
}
