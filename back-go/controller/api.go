package controller

import (
	"gin-web/models/authcCenter"
	"gin-web/utils"
	"gin-web/utils/extendController"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ApiController struct {
	extendController.BaseController
}

func (a *ApiController) Insert(c *gin.Context) {
	//参数校验
	var api authcCenter.Api
	if err := c.Bind(&api); err != nil {
		a.SendParameterErrorResponse(c, 4002, err)
		return
	}
	if api.Method != "POST" && api.Method != "GET" && api.Method != "DELETE" && api.Method != "PUT" {
		a.SendParameterErrorResponse(c, 4003, nil)
		return
	}
	//DB操作
	isExist, err := api.IsExist()
	if isExist || err != nil {
		a.SendServerErrorResponse(c, 5101, err)
		return
	}
	if err = api.Insert(); err != nil {
		a.SendServerErrorResponse(c, 5100, err)
		return
	}
	a.SendSuccessResponse(c, "success")
}

func (a *ApiController) Remove(c *gin.Context) {
	//参数校验
	id := c.Query("id")
	if id == "" {
		a.SendParameterErrorResponse(c, 4001, nil)
		return
	}
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		a.SendCustomResponseByBacked(c, "参数格式转化错误", "Parameter format conversion error", err)
		return
	}
	//DB操作
	if err := new(authcCenter.Api).Remove(idInt64); err != nil {
		a.SendServerErrorResponse(c, 5110, err)
		return
	}
	a.SendSuccessResponse(c, "success")
}

func (a *ApiController) Edit(c *gin.Context) {
	//参数校验
	var api authcCenter.Api
	if err := c.Bind(&api); err != nil {
		a.SendParameterErrorResponse(c, 4002, err)
		return
	}
	if api.Method != "POST" && api.Method != "GET" && api.Method != "DELETE" && api.Method != "PUT" {
		a.SendParameterErrorResponse(c, 4003, nil)
		return
	}
	//DB操作
	if err := api.Edit(); err != nil {
		a.SendServerErrorResponse(c, 5120, err)
		return
	}
	a.SendSuccessResponse(c, "success")
}

// 查询
func (a *ApiController) GetList(c *gin.Context) {
	//参数校验
	var api authcCenter.Api
	api.Name, api.Url = c.Query("name"), c.Query("url")
	currPage, pageSize := c.DefaultQuery("currPage", "1"), c.DefaultQuery("pageSize", "10")
	startTime, endTime := c.Query("startTime"), c.Query("endTime")
	skip, limit, err := utils.GetPage(currPage, pageSize)
	if err != nil {
		a.SendCustomResponseByBacked(c, "分页失败", "Paging failed", err)
		return
	}
	//DB操作
	resDB, count, err := api.GetList(skip, limit, startTime, endTime)
	if err != nil {
		a.SendServerErrorResponse(c, 5130, err)
		return
	}
	a.SendSuccessResponse(c, struct {
		Logs  []authcCenter.Api `json:"apis"`
		Total int64             `json:"total"`
	}{
		resDB,
		count,
	})
}
