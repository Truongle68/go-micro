package v1

import (
	"github.com/TruongLe68/go-micro/pkg/ginmw"
	"github.com/gin-gonic/gin"
)

func (r *V1) getUserID(c *gin.Context) (string, bool) {
	return ginmw.GetUserID(c)
}

func (r *V1) getRole(c *gin.Context) (string, bool) {
	return ginmw.GetRole(c)
}
