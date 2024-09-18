package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) router() *gin.Engine {
	r := gin.Default()

	//v1
	v1 := r.Group("/api/v1")
	{
		v1.GET("/", s.index)

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
