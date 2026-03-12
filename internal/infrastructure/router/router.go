package router

import (
	"project-bcc/internal/adapter/handler"
	"project-bcc/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Router struct {
	authHandler *handler.AuthHandler
}

func NewRouter(authHandler *handler.AuthHandler) *Router {
	return &Router{authHandler: authHandler}
}

func (r *Router) SetupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())

	auth := router.Group("/api/auth")
	{
		auth.POST("/register", r.authHandler.Register)
		auth.POST("/login", r.authHandler.Login)
	}

	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{

	}
	return router
}
