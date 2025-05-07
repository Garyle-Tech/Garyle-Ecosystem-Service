package middleware

import (
	"net/http"

	"ecosystem.garyle/service/pkg/utils/response"
	"github.com/gin-gonic/gin"
)

func ResponseMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if c.Writer.Written() || c.IsAborted() {
			return
		}

		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			code := http.StatusInternalServerError
			if c.Writer.Status() != http.StatusOK {
				code = c.Writer.Status()
			}

			resp := response.NewErrorResponse(code, err.Error())
			c.JSON(code, resp)
			return
		}

		data, exists := c.Get("data")
		if !exists {
			data = gin.H{}
		}

		message, exists := c.Get("message")
		msg := ""
		if exists {
			msg = message.(string)
		}

		page, pageExists := c.Get("page")
		limit, limitExists := c.Get("limit")
		total, totalExists := c.Get("total")
		lastPage, lastPageExists := c.Get("last_page")

		if pageExists && limitExists && totalExists && lastPageExists {
			pagination := response.Pagination{
				Page:     page.(int),
				Limit:    limit.(int),
				Total:    total.(int),
				LastPage: lastPage.(int),
			}
			resp := response.NewSuccessResponseWithPagination(data, msg, pagination)
			c.JSON(http.StatusOK, resp)
		} else {
			resp := response.NewSuccessResponse(data, msg)
			c.JSON(http.StatusOK, resp)
		}
	}
}
