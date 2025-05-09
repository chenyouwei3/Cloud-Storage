package controller

import (
	"gin-web/models/authcCenter"
	"gin-web/utils"
	"gin-web/utils/extendController"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoleController struct {
	extendController.BaseController
}

type roleRequest struct {
	Role        authcCenter.Role
	AddApis     []int `json:"addApis"`
	DeletedApis []int `json:"deletedApis"`
}

func (r *RoleController) Insert(c *gin.Context) {
	//参数校验
	var role roleRequest
	if err := c.Bind(&role); err != nil {
		r.SendParameterErrorResponse(c, 4002, err)
		return
	}
	//DB操作
	isExist, err := role.Role.IsExist()
	if isExist || err != nil {
		r.SendServerErrorResponse(c, 5101, err)
		return
	}
	if err = role.Role.Insert(role.AddApis); err != nil {
		r.SendServerErrorResponse(c, 5100, err)
		return
	}
	r.SendSuccessResponse(c, "success")
}

func (r *RoleController) Remove(c *gin.Context) {
	//参数校验
	id := c.Query("id")
	if id == "" {
		r.SendParameterErrorResponse(c, 4001, nil)
		return
	}
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		r.SendCustomResponseByBacked(c, "参数格式转化错误", "Parameter format conversion error", err)
		return
	}
	//DB操作
	if err := new(authcCenter.Role).Remove(idInt64); err != nil {
		r.SendServerErrorResponse(c, 5110, err)
		return
	}
	r.SendSuccessResponse(c, "success")
}

func (r *RoleController) Edit(c *gin.Context) {
	//参数校验
	var role roleRequest
	if err := c.Bind(&role); err != nil {
		r.SendParameterErrorResponse(c, 4002, err)
		return
	}
	//DB操作
	if err := role.Role.Edit(role.AddApis, role.DeletedApis); err != nil {
		r.SendServerErrorResponse(c, 5120, err)
		return
	}
	r.SendSuccessResponse(c, "success")
}

// 查询
func (r *RoleController) GetList(c *gin.Context) {
	//参数校验
	var role authcCenter.Role
	role.Name = c.Query("name") //角色名称
	currPage := c.DefaultQuery("currPage", "1")
	pageSize := c.DefaultQuery("pageSize", "10")
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")
	skip, limit, err := utils.GetPage(currPage, pageSize)
	if err != nil {
		r.SendCustomResponseByBacked(c, "分页失败", "Paging failed", err)
		return
	}
	//DB操作
	resDB, count, err := role.GetList(skip, limit, startTime, endTime)
	if err != nil {
		r.SendServerErrorResponse(c, 5130, err)
		return
	}
	r.SendSuccessResponse(c, struct {
		Logs  []authcCenter.Role `json:"roles"`
		Total int64              `json:"total"`
	}{
		resDB,
		count,
	})
}
