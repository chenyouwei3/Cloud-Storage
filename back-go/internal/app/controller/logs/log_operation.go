package logs

import (
	"gin-web/internal/app/models"
	"gin-web/pkg"
	"gin-web/pkg/extendController"
	"github.com/gin-gonic/gin"
)

type OperationLogController struct {
	extendController.BaseController
}

func (o *OperationLogController) GetList(c *gin.Context) {
	//参数校验
	var log models.OperationLog
	log.Account = c.Query("account")
	currPage, pageSize := c.DefaultQuery("currPage", "1"), c.DefaultQuery("pageSize", "10")
	startTime, endTime := c.Query("startTime"), c.Query("endTime")
	skip, limit, err := pkg.GetPage(currPage, pageSize)
	if err != nil {
		o.SendCustomResponseByBacked(c, "分页失败", "Paging failed", err)
		return
	}
	//DB操作
	resDB, count, err := log.GetList(skip, limit, startTime, endTime)
	if err != nil {
		o.SendServerErrorResponse(c, 5130, err)
		return
	}
	o.SendSuccessResponse(c, struct {
		Logs  []models.OperationLog `json:"logs"`
		Total int64                 `json:"total"`
	}{
		resDB,
		count,
	})
}
