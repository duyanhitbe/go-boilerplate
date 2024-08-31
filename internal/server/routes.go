package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() *gin.Engine {
	r := gin.Default()

	//v1
	v1 := r.Group("/api/v1")
	{
		v1.GET("/", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "Todo app",
				"version": "1.0.0",
				"status":  "ok",
			})
		})

		//todo endpoints
		todo := r.Group("/todo")
		{
			todo.GET("/", s.GetAllTodo)
		}
	}

	return r
}
