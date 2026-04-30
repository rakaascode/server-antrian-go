package crm

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CrmHandler struct {
	service CrmService
}

func NewCrmHandler(s CrmService) *CrmHandler {
	return &CrmHandler{s}
}

// ManualSend POST /api/crm/send
// Admin input no WA dan pesan secara manual (tidak terkait data antrian)
func (h *CrmHandler) ManualSend(c *gin.Context) {
	var req ManualSendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	if err := h.service.ManualSend(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Pesan berhasil dikirim ke " + req.NoWA,
	})
}

// ReminderFromAntrian POST /api/crm/reminders
// Kirim pengingat ke pelanggan berdasarkan data antrian (ambil no WA dari antrian)
func (h *CrmHandler) ReminderFromAntrian(c *gin.Context) {
	var req ReminderFromAntrianRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	cabangID := c.MustGet("cabang_id").(uint)
	if err := h.service.ReminderFromAntrian(req, cabangID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Pengingat berhasil dikirim ke pelanggan",
	})
}
