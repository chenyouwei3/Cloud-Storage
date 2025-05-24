package controller

import (
	"fmt"
	"gin-web/internal/app/models/authCenter"
	mysqlDB "gin-web/internal/initialize/mysql"
	"gin-web/pkg"
	"gin-web/pkg/extendController"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
)

type RoleController struct {
	extendController.BaseController
}

type roleRequest struct {
	Role        authCenter.Role `json:"role"`
	AddApis     []int           `json:"addApis"`
	DeletedApis []int           `json:"deletedApis"`
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
	var req struct {
		Id int64 `json:"id"`
	}
	if err := c.Bind(&req); err != nil {
		r.SendParameterErrorResponse(c, 4002, err)
		return
	}
	//DB操作
	if err := new(authCenter.Role).Remove(req.Id); err != nil {
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
	fmt.Println("TESTING", role)
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
	var role authCenter.Role
	role.Name = c.Query("name") //角色名称
	currPage := c.DefaultQuery("currPage", "1")
	pageSize := c.DefaultQuery("pageSize", "10")
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")
	skip, limit, err := pkg.GetPage(currPage, pageSize)
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
		Logs  []authCenter.Role `json:"roles"`
		Total int64             `json:"total"`
	}{
		resDB,
		count,
	})
}

func (r *RoleController) GetApisByRoleID(c *gin.Context) {
	roleId := c.Query("id")
	if roleId == "" {
		r.SendParameterErrorResponse(c, 4001, nil)
		return
	}
	id, err := strconv.Atoi(roleId)
	if err != nil {
		r.SendParameterErrorResponse(c, 4003, err)
		return
	}
	var role authCenter.Role
	err = mysqlDB.DB.
		Preload("Apis", func(db *gorm.DB) *gorm.DB {
			return db.Select("api.id", "api.name") // 指定表名.字段更稳妥
		}).
		Where("id = ?", id).
		First(&role).Error
	if err != nil {
		r.SendServerErrorResponse(c, 5130, err)
		return
	}
	r.SendSuccessResponse(c, role.Apis)
}
