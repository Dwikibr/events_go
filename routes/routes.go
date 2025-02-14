package routes

import (
	"RestApi/middleware"
	"github.com/gin-gonic/gin"
)

func SetRoutes(server *gin.Engine) {
	authenticated := server.Group("/")
	authenticated.Use(middleware.AuthenticateUser)
	authenticated.POST("/events", createEvent)
	authenticated.PUT("/events/:id", updateEvent)
	authenticated.DELETE("/events/:id", deleteEvent)
	authenticated.POST("/events/:id/register", RegisterEvent)
	authenticated.DELETE("/events/:id/cancel-register", CancelRegistration)
	authenticated.GET("/event/:id/detail", GetEventDetail)

	server.GET("/events", getEvents)
	server.GET("/events/:id", getEvent)
	server.POST("/signup", SignUp)
	server.POST("/login", Login)
}
