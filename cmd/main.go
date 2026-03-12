package main

import (
	"log"
	"project-bcc/internal/adapter/handler"
	"project-bcc/internal/adapter/repository"
	"project-bcc/internal/infrastructure/database"
	"project-bcc/internal/infrastructure/router"
	"project-bcc/internal/usecase"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Gagal load file env:", err)
	}

	db := database.ConnectDB()
	authRepo := repository.NewAuthRepository(db)
	authUsecase := usecase.NewAuthUsecase(authRepo)
	authHandler := handler.NewAuthHandler(authUsecase)
	r := router.NewRouter(authHandler)
	app := r.SetupRouter()

	log.Println("Server berjalan di port 8000")
	err = app.Run(":8000")
	if err != nil {
		log.Fatal("Gagal menjalankan server:", err)
	}

}
