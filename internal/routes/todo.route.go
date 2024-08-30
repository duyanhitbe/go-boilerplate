package routes

import (
	"github.com/duyanhitbe/go-boilerplate/internal/handlers"
	"github.com/gin-gonic/gin"
)

func (s *Server) registerTodoRoutes(r *gin.Engine) {
	//init TodoHandler
	handler := handlers.NewTodoHandler(s.store)

	//v1 group
	v1 := r.Group("/api/v1")
	{
		//todo group
		todo := v1.Group("/todo")
		{
			todo.GET("/", handler.GetAllTodo)
		}
	}
}
