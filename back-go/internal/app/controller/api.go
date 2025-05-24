package controller

import (
	"gin-web/internal/app/models/authCenter"
	"gin-web/pkg"
	"gin-web/pkg/extendController"
	"github.com/gin-gonic/gin"
)

type ApiController struct {
	extendController.BaseController
}

func (a *ApiController) Insert(c *gin.Context) {
	//参数校验
	var api authCenter.Api
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
	var req struct {
		Id int64 `json:"id"`
	}
	if err := c.Bind(&req); err != nil {
		a.SendParameterErrorResponse(c, 4002, err)
		return
	}
	//DB操作
	if err := new(authCenter.Api).Remove(req.Id); err != nil {
		a.SendServerErrorResponse(c, 5110, err)
		return
	}
	a.SendSuccessResponse(c, "success")
}

func (a *ApiController) Edit(c *gin.Context) {
	//参数校验
	var api authCenter.Api
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
	var api authCenter.Api
	api.Name, api.Url, api.Method = c.Query("name"), c.Query("url"), c.Query("method")
	currPage, pageSize := c.DefaultQuery("currPage", "1"), c.DefaultQuery("pageSize", "10")
	startTime, endTime := c.Query("startTime"), c.Query("endTime")
	skip, limit, err := pkg.GetPage(currPage, pageSize)
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
		Logs  []authCenter.Api `json:"apis"`
		Total int64            `json:"total"`
	}{
		resDB,
		count,
	})
}
