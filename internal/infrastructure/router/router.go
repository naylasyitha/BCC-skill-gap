package router

import (
	"project-bcc/internal/adapter/handler"
	"project-bcc/internal/middleware"
	"project-bcc/internal/usecase"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Router struct {
	authHandler          *handler.AuthHandler
	careerHandler        *handler.CareerHandler
	skillHandler         *handler.SkillHandler
	selfAssessment       *handler.SelfAssessmentHandler
	quizHandler          *handler.QuizHandler
	careerSessionHandler *handler.CareerSessionHandler
	questionHandler      *handler.QuestionHandler
	authRepo             usecase.AuthRepository
}

func NewRouter(
	ah *handler.AuthHandler,
	ch *handler.CareerHandler,
	sh *handler.SkillHandler,
	sa *handler.SelfAssessmentHandler,
	qu *handler.QuizHandler,
	cs *handler.CareerSessionHandler,
	que *handler.QuestionHandler,
	ar usecase.AuthRepository,
) *Router {
	return &Router{
		authHandler:          ah,
		careerHandler:        ch,
		skillHandler:         sh,
		selfAssessment:       sa,
		quizHandler:          qu,
		careerSessionHandler: cs,
		questionHandler:      que,
		authRepo:             ar,
	}
}

func (r *Router) SetupRouter() *gin.Engine {
	router := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = false
	corsConfig.AllowOrigins = []string{
		"http://localhost:3000",
	}
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{
		"Origin",
		"Content-Type",
		"Accept",
		"Authorization",
	}

	router.Use(cors.New(corsConfig))

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

	api.Use(middleware.AuthMiddleware(r.authRepo))

	career := api.Group("/careers")
	{
		career.GET("", r.careerHandler.GetAllCareer)
		career.GET("/:careerId", r.careerHandler.GetCareerById)

		career.POST("", middleware.AdminMiddleware(), r.careerHandler.CreateCareer)
		career.PATCH("/:careerId", middleware.AdminMiddleware(), r.careerHandler.UpdateCareer)
		career.DELETE("/:careerId", middleware.AdminMiddleware(), r.careerHandler.DeleteCareer)
	}

	skill := api.Group("/skills")
	{
		skill.GET("", r.skillHandler.GetAllSkill)
		skill.GET("/:skillId", r.skillHandler.GetSkillById)

		skill.POST("", middleware.AdminMiddleware(), r.skillHandler.CreateSkill)
		skill.PATCH("/:skillId", middleware.AdminMiddleware(), r.skillHandler.UpdateSkill)
		skill.DELETE("/:skillId", middleware.AdminMiddleware(), r.skillHandler.DeleteSkill)

	}

	question := api.Group("/questions")
	{
		question.GET("", r.questionHandler.GetAllQuestion)
		question.GET("/:questionId", r.questionHandler.GetQuestionById)

		question.POST("", middleware.AdminMiddleware(), r.questionHandler.CreateQuestion)
		question.PATCH("/:questionId", middleware.AdminMiddleware(), r.questionHandler.UpdateQuestion)
		question.DELETE("/:questionId", middleware.AdminMiddleware(), r.questionHandler.DeleteQuestion)

	}

	careerSession := api.Group("/career-sessions")
	{
		careerSession.POST("", r.careerSessionHandler.Create)
		careerSession.GET("/:careerSessionId", r.careerSessionHandler.GetCareerSession)
		careerSession.POST("/:careerSessionId/assessment", r.selfAssessment.SubmitAssessment)
		careerSession.POST("/quiz/:careerSessionId/start", r.quizHandler.StartQuiz)
	}

	//konsep awal untuk cadangan just in case
	// careerSkills := api.Group("/career-skills")
	// {
	// 	careerSkills.POST("", middleware.AdminMiddleware(), r.skillHandler.CareerSkillAsign)
	// 	careerSkills.PATCH("/:id", middleware.AdminMiddleware(), r.skillHandler.UpdateCareerSkill)
	// 	careerSkills.DELETE("/:id", middleware.AdminMiddleware(), r.skillHandler.RemoveSkillFromCareer)
	// }

	return router
}
