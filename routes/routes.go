package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rakaascode/server-antrian-go.git/handler"
)


func SetupRoutes(r *gin.Engine, h *handler.UserHandler) {
	api := r.Group("/api")

	api.GET("/users", h.GetAll)
	api.GET("/users/:id", h.GetByID)
	api.POST("/users", h.Create)
	api.PUT("/users/:id", h.Update)
	api.DELETE("/users/:id", h.Delete)
}