package controller

import (
	"gin-web/utils/extendController"
	"github.com/gin-gonic/gin"
)

type DefaultController struct {
	extendController.BaseController
}

// HandleNotFound 405处理
func (d DefaultController) HandleNotFound(c *gin.Context) {
	d.SendMethodNotAllowedResponse(c)
}
