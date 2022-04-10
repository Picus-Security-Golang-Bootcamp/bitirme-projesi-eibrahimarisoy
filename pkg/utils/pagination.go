package utils

import (
	"fmt"
	"patika-ecommerce/internal/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GeneratePaginationFromRequest(c *gin.Context) model.Pagination {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	fmt.Println("page", page)
	if page == 0 {
		page = 1
	}
	sort := c.DefaultQuery("sort", "created_at desc")
	q := c.DefaultQuery("q", "")

	return model.Pagination{
		Limit: limit,
		Page:  page,
		Sort:  sort,
		Q:     q,
	}

}
