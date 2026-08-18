package v1

import (
	"github.com/TruongLe68/go-micro/pkg/ginmw"
	"github.com/gin-gonic/gin"
)

func (r *V1) getUserID(c *gin.Context) (string, bool) {
	return ginmw.GetUserID(c)
}
