package router

import (
	"project-bcc/internal/adapter/handler"
	"project-bcc/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Router struct {
	authHandler    *handler.AuthHandler
	careerHandler  *handler.CareerHandler
	skillHandler   *handler.SkillHandler
	selfAssessment *handler.SelfAssessmentHandler
	quizHandler    *handler.QuizHandler
}

func NewRouter(
	ah *handler.AuthHandler,
	ch *handler.CareerHandler,
	sh *handler.SkillHandler,
	sa *handler.SelfAssessmentHandler,
	qu *handler.QuizHandler,
) *Router {
	return &Router{
		authHandler:    ah,
		careerHandler:  ch,
		skillHandler:   sh,
		selfAssessment: sa,
		quizHandler:    qu,
	}
}

func (r *Router) SetupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())

	api := router.Group("/api")

	auth := router.Group("/api/auth")
	{
		auth.POST("/register", r.authHandler.Register)
		auth.POST("/login", r.authHandler.Login)
		auth.POST("/verify", r.authHandler.VerifyEmail)
		auth.POST("/resend-verify", r.authHandler.ResendVerification)
		auth.POST("/refresh", r.authHandler.Refresh)
		auth.POST("/logout", r.authHandler.Logout)
		auth.POST("/forgot-password", r.authHandler.ForgotPassword)
		auth.POST("/reset-password", r.authHandler.ResetPassword)
	}

	api.Use(middleware.AuthMiddleware())

	career := api.Group("/careers")
	{
		career.GET("", r.careerHandler.GetAllCareer)
		career.GET("/:id", r.careerHandler.GetCareerById)

		career.POST("", middleware.AdminMiddleware(), r.careerHandler.CreateCareer)
		career.PATCH("/:id", middleware.AdminMiddleware(), r.careerHandler.UpdateCareer)
		career.DELETE("/:id", middleware.AdminMiddleware(), r.careerHandler.DeleteCareer)
	}

	skill := api.Group("/skills")
	{
		skill.GET("", r.skillHandler.GetAllSkill)
		skill.GET("/:id", r.skillHandler.GetSkillById)

		skill.POST("", middleware.AdminMiddleware(), r.skillHandler.CreateSkill)
		skill.PATCH("/:id", middleware.AdminMiddleware(), r.skillHandler.UpdateSkill)
		skill.DELETE("/:id", middleware.AdminMiddleware(), r.skillHandler.DeleteSkill)

	}

	careerSession := api.Group("/career-sessions")
	{
		careerSession.POST("/:id/assessment", r.selfAssessment.SubmitAssessment)
		careerSession.POST("/:id/quiz/start", r.quizHandler.StartQuiz)
	}

	careerSkills := api.Group("/career-skills")
	{
		careerSkills.POST("/career-skills", middleware.AdminMiddleware(), r.skillHandler.CareerSkillAsign)
		careerSkills.PATCH("/career-skills/:id", middleware.AdminMiddleware(), r.skillHandler.UpdateCareerSkill)
		careerSkills.DELETE("/career-skills/:id", middleware.AdminMiddleware(), r.skillHandler.RemoveSkillFromCareer)
	}

	return router
}
