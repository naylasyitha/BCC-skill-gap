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
	careerRepo := repository.NewCareerRepository(db)
	skillRepo := repository.NewSkillRepository(db)
	careerSkillRepo := repository.NewCareerSkillRepository(db)
	selfAssesmentRepo := repository.NewSelfAssessmentRepository(db)
	quizRepo := repository.NewQuizRepository(db)
	careerSessionRepo := repository.NewCareerSessionRepository(db)
	questionRepo := repository.NewQuestionRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)

	authUsecase := usecase.NewAuthUsecase(authRepo)
	careerUsecase := usecase.NewCareerUsecase(careerRepo, skillRepo)
	skillUsecase := usecase.NewSkillUsecase(skillRepo, careerSkillRepo)
	selfAssessmentUsecase := usecase.NewSelfAssessmentUsecase(selfAssesmentRepo, careerSessionRepo)
	quizUsecase := usecase.NewQuizUsecase(quizRepo)
	careerSessionUsecase := usecase.NewCareerSessionUsecase(careerSessionRepo, careerRepo, authRepo)
	questionUsecase := usecase.NewQuestionUsecase(questionRepo, skillRepo)
	userUsecase := usecase.NewUserUsecase(authRepo)
	paymentUsecase := usecase.NewPaymentUsecase(paymentRepo)

	authHandler := handler.NewAuthHandler(authUsecase)
	careerHandler := handler.NewCareerHandler(careerUsecase)
	skillHandler := handler.NewSkillHandler(skillUsecase)
	selfAssesmentHandler := handler.NewSelfAssessmentHandler(selfAssessmentUsecase)
	quizHandler := handler.NewQuizHandler(quizUsecase)
	careerSessionHandler := handler.NewCareerSessionHandler(careerSessionUsecase)
	questionHandler := handler.NewQuestionHandler(questionUsecase)
	userHandler := handler.NewUserHandler(userUsecase)
	paymentHandler := handler.NewPaymentHandler(paymentUsecase)

	r := router.NewRouter(authHandler, careerHandler, skillHandler, selfAssesmentHandler, quizHandler, careerSessionHandler, questionHandler, userHandler, paymentHandler, authRepo)

	app := r.SetupRouter()

	log.Println("Server berjalan di port 8000")
	err = app.Run(":8000")
	if err != nil {
		log.Fatal("Gagal menjalankan server:", err)
	}

}
