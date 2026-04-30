package broadcast

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BroadcastHandler struct {
	service BroadcastService
}

func NewBroadcastHandler(s BroadcastService) *BroadcastHandler {
	return &BroadcastHandler{s}
}

// Create POST /api/broadcast  (admin only)
func (h *BroadcastHandler) Create(c *gin.Context) {
	var req CreateBroadcastRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	adminID := c.MustGet("user_id").(uint)
	b, err := h.service.Create(req, adminID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Broadcast berhasil dikirim",
		"data":    b,
	})
}

// GetAll GET /api/broadcast/all  (admin only)
func (h *BroadcastHandler) GetAll(c *gin.Context) {
	list, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": list})
}

// GetForUser GET /api/broadcast  (user login)
// Promo global + antrian dari cabang yang pernah dipakai user
func (h *BroadcastHandler) GetForUser(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	list, err := h.service.GetForUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": list})
}

// GetByID GET /api/broadcast/:id  (user / admin)
// Detail lengkap satu broadcast
func (h *BroadcastHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "id tidak valid"})
		return
	}
	detail, err := h.service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": detail})
}
