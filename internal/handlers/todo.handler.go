package handlers

import (
	db "github.com/duyanhitbe/go-boilerplate/internal/database/generated"
	"github.com/gin-gonic/gin"
)

type TodoHandler struct {
	store db.Store
}

func NewTodoHandler(store db.Store) *TodoHandler {
	return &TodoHandler{
		store: store,
	}
}

func (h *TodoHandler) GetAllTodo(ctx *gin.Context) {
	data, err := h.store.GetAllTodo(ctx)
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
