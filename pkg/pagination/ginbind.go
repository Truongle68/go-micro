package pagination

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func FromQuery(c *gin.Context) Params {
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
	return Params{
		Page:  page,
		Limit: limit,
	}.Normalize()
}
