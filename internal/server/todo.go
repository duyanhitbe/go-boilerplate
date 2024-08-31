package server

import "github.com/gin-gonic/gin"

// /api/v1/todo
func (s *Server) GetAllTodo(ctx *gin.Context) {
	data, err := s.store.GetAllTodo(ctx)
	if err != nil {
		ctx.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"status": 200,
		"code":   "ok",
		"data":   data,
	})
}
