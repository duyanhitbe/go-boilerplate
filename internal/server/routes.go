package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) router() *gin.Engine {
	r := gin.Default()
	r.GET("/", s.index)

	//v1 endpoints
	v1 := r.Group("/api/v1")
	{
		//auth endpoints
		auth := v1.Group("/auth")
		{
			auth.POST("/register", s.register)
			auth.POST("/login", s.login)
			auth.GET("/me", s.authenticationMiddleware, s.me)
		}
	}

	return r
}

func (s *Server) index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Go boilerplate",
		"version": "1.0.0",
		"status":  "ok",
	})
}
