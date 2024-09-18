package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, response{
		Status:  http.StatusCreated,
		Message: http.StatusText(http.StatusCreated),
		Data:    data,
	})
}

func ok(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, response{
		Status:  http.StatusOK,
		Message: http.StatusText(http.StatusOK),
		Data:    data,
	})
}
