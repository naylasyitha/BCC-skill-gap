package router

import (
	"project-bcc/internal/adapter/handler"
	"project-bcc/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Router struct {
	authHandler   *handler.AuthHandler
	careerHandler *handler.CareerHandler
	skillHandler  *handler.SkillHandler
}

func NewRouter(
	ah *handler.AuthHandler,
	ch *handler.CareerHandler,
	sh *handler.SkillHandler,
) *Router {
	return &Router{
		authHandler:   ah,
		careerHandler: ch,
		skillHandler:  sh,
	}
}

func (r *Router) SetupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())

	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{

	}

	auth := router.Group("/api/auth")
	{
		auth.POST("/register", r.authHandler.Register)
		auth.POST("/login", r.authHandler.Login)
	}

	career := router.Group("/api/careers")
	{
		career.GET("", r.careerHandler.GetAllCareer)
		career.GET("/:id/skills", r.careerHandler.GetCareerById)
	}

	skill := router.Group("/api/skills")
	{
		skill.GET("", r.skillHandler.GetAllSkill)
		skill.GET("/:id", r.skillHandler.GetSkillById)
	}

	admin := api.Group("/admin")
	admin.Use(middleware.AdminMiddleware())
	{
		admin.POST("/careers", r.careerHandler.CreateCareer)
		admin.PATCH("/careers/:id", r.careerHandler.UpdateCareer)
		admin.DELETE("/careers/:id", r.careerHandler.DeleteCareer)

		admin.POST("/skills", r.skillHandler.CreateSkill)
		admin.PATCH("/skills/:id", r.skillHandler.UpdateSkill)
		admin.DELETE("/skills/:id", r.skillHandler.DeleteSkill)

		admin.POST("/career-skills", r.skillHandler.CareerSkillAsign)
		admin.PATCH("/career-skills/:id", r.skillHandler.UpdateCareerSkill)
		admin.DELETE("/career-skills/:id", r.skillHandler.RemoveSkillFromCareer)
	}

	return router
}
