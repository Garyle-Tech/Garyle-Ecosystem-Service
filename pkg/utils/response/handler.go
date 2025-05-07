package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, data interface{}, message string) {
	resp := NewSuccessResponse(data, message)
	c.JSON(http.StatusOK, resp)
}

func SuccessWithPagination(c *gin.Context, data interface{}, message string, page, limit, total int) {
	lastPage := total / limit
	if total%limit > 0 {
		lastPage++
	}

	pagination := Pagination{
		Page:     page,
		Limit:    limit,
		Total:    total,
		LastPage: lastPage,
	}

	resp := NewSuccessResponseWithPagination(data, message, pagination)
	c.JSON(http.StatusOK, resp)
}

func Created(c *gin.Context, data interface{}, message string) {
	resp := Response{
		Meta: Meta{
			Message: message,
			Status:  "success",
			Code:    http.StatusCreated,
		},
		Data: data,
	}
	c.JSON(http.StatusCreated, resp)
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func Error(c *gin.Context, code int, message string) {
	resp := NewErrorResponse(code, message)
	c.JSON(code, resp)
}

func BadRequest(c *gin.Context, message string) {
	resp := BadRequestError(message)
	c.JSON(http.StatusBadRequest, resp)
}

func Unauthorized(c *gin.Context, message string) {
	resp := UnauthorizedError(message)
	c.JSON(http.StatusUnauthorized, resp)
}

func Forbidden(c *gin.Context, message string) {
	resp := ForbiddenError(message)
	c.JSON(http.StatusForbidden, resp)
}

func NotFound(c *gin.Context, message string) {
	resp := NotFoundError(message)
	c.JSON(http.StatusNotFound, resp)
}

func Server(c *gin.Context, message string) {
	resp := ServerError(message)
	c.JSON(http.StatusInternalServerError, resp)
}
