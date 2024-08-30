package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) registerIndexRoutes(r *gin.Engine) {
	//v1 group
	v1 := r.Group("/api/v1")
	{
		v1.GET("/", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "Todo app",
				"version": "1.0.0",
				"status":  "ok",
			})
		})
	}
}
