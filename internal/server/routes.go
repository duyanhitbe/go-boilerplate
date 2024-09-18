package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) router() *gin.Engine {
	r := gin.Default()

	//v1
	v1 := r.Group("/api/v1")
	{
		v1.GET("/", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "Go boilerplate",
				"version": "1.0.0",
				"status":  "ok",
			})
		})

		//auth endpoints
		auth := v1.Group("/auth")
		{
			auth.POST("/register", s.register)
			auth.POST("/login", s.login)
		}
	}

	return r
}
