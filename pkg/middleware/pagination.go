package mw

import (
	"patika-ecommerce/pkg/pagination"
	"strconv"

	"github.com/gin-gonic/gin"
)

// PaginationMiddleware is a middleware to paginate the results
func PaginationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		pagination := &pagination.Pagination{}

		limit, err := strconv.Atoi(c.Query("limit"))
		if err != nil {
			limit = 10
		}
		pagination.Limit = limit
		page, err := strconv.Atoi(c.Query("page"))
		if err != nil {
			page = 0
		}
		pagination.Page = page

		sort := c.Query("sort")
		if sort == "" {
			sort = "created_at desc"
		}
		pagination.Sort = sort
		pagination.Q = c.DefaultQuery("q", "")

		c.Set("pagination", pagination)

		c.Next()
		return
	}
}
