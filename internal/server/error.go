package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type errorResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Error   string      `json:"error,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

func throwBadRequest(c *gin.Context, err error) {
	c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{
		Status:  http.StatusBadRequest,
		Message: http.StatusText(http.StatusBadRequest),
		Error:   err.Error(),
	})
}

func throwValidationError(c *gin.Context, errs interface{}) {
	c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{
		Status:  http.StatusBadRequest,
		Message: http.StatusText(http.StatusBadRequest),
		Errors:  errs,
	})
}

func throwForbidden(c *gin.Context, err error) {
	c.AbortWithStatusJSON(http.StatusForbidden, errorResponse{
		Status:  http.StatusForbidden,
		Message: http.StatusText(http.StatusForbidden),
		Error:   err.Error(),
	})
}

func throwConflict(c *gin.Context, err error) {
	c.AbortWithStatusJSON(http.StatusConflict, errorResponse{
		Status:  http.StatusConflict,
		Message: http.StatusText(http.StatusConflict),
		Error:   err.Error(),
	})
}

func throwInternalServerError(c *gin.Context, err error) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse{
		Status:  http.StatusInternalServerError,
		Message: http.StatusText(http.StatusInternalServerError),
		Error:   err.Error(),
	})
}
