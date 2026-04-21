package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rakaascode/server-antrian-go.git/config"
	"github.com/rakaascode/server-antrian-go.git/handler"
	"github.com/rakaascode/server-antrian-go.git/models"
	"github.com/rakaascode/server-antrian-go.git/repository"
	"github.com/rakaascode/server-antrian-go.git/routes"
	"github.com/rakaascode/server-antrian-go.git/service"
)

func main(){

	godotenv.Load()

	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal("Gagal konek DB:", err)
	}

	log.Println("DB connected...")

	// migrate
	db.AutoMigrate(&models.User{})

	log.Println("Migrasi selesai...")


	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	route := gin.Default()

	routes.SetupRoutes(route, userHandler)

	route.Run(":8080")
}

