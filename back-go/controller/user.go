package controller

import (
	"gin-web/middleware"
	"gin-web/models/authcCenter"
	"gin-web/utils"
	"gin-web/utils/extendController"
	"strconv"

	"github.com/gin-gonic/gin"
)

type userRequest struct {
	User         authcCenter.User `json:"user"`
	AddRoles     []int            `json:"addRoles"`
	DeletedRoles []int            `json:"deletedRoles"`
}

type UserController struct {
	extendController.BaseController
}

func (u *UserController) Insert(c *gin.Context) {
	//参数校验
	var reqUser userRequest
	if err := c.BindJSON(&reqUser); err != nil {
		u.SendParameterErrorResponse(c, 4002, err)
		return
	}
	//DB操作
	isExist, err := reqUser.User.IsExist()
	if isExist || err != nil {
		u.SendServerErrorResponse(c, 5101, err)
		return
	}
	if err = reqUser.User.Insert(reqUser.AddRoles); err != nil {
		u.SendServerErrorResponse(c, 5100, err)
		return
	}
	u.SendSuccessResponse(c, "success")
}

func (u *UserController) Remove(c *gin.Context) {
	//参数校验
	id := c.Query("id")
	if id == "" {
		u.SendParameterErrorResponse(c, 4001, nil)
		return
	}
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		u.SendCustomResponseByBacked(c, "参数格式转化错误", "Parameter format conversion error", err)
		return
	}
	//DB操作
	if err = new(authcCenter.User).Remove(idInt64); err != nil {
		u.SendServerErrorResponse(c, 5110, err)
		return
	}
	u.SendSuccessResponse(c, "success")
}

func (u *UserController) Edit(c *gin.Context) {
	//参数校验
	var reqUser userRequest
	if err := c.ShouldBind(&reqUser); err != nil {
		u.SendParameterErrorResponse(c, 4002, err)
		return
	}
	//DB操作
	if err := reqUser.User.Edit(reqUser.AddRoles, reqUser.DeletedRoles); err != nil {
		u.SendServerErrorResponse(c, 5120, err)
		return
	}
	u.SendSuccessResponse(c, "success")
}

func (u *UserController) GetList(c *gin.Context) {
	//参数校验
	var user authcCenter.User
	user.Name = c.Query("name")
	currPage, pageSize := c.DefaultQuery("currPage", "1"), c.DefaultQuery("pageSize", "10")
	startTime, endTime := c.Query("startTime"), c.Query("endTime")
	skip, limit, err := utils.GetPage(currPage, pageSize)
	if err != nil {
		u.SendCustomResponseByBacked(c, "分页失败", "Paging failed", err)
		return
	}
	//DB操作
	resDB, count, err := user.GetList(skip, limit, startTime, endTime)
	if err != nil {
		u.SendServerErrorResponse(c, 5130, err)
		return
	}
	u.SendSuccessResponse(c, struct {
		Logs  []authcCenter.User `json:"users"`
		Total int64              `json:"total"`
	}{
		resDB,
		count,
	})
}

func (u *UserController) Login(c *gin.Context) {
	//参数校验
	var reqUser struct {
		Account  string `form:"account"`
		Password string `form:"password"`
	}
	if err := c.ShouldBind(&reqUser); err != nil {
		u.SendParameterErrorResponse(c, 4002, err)
		return
	}
	user, err := new(authcCenter.User).GetOne(reqUser.Account, "")
	if err != nil {
		u.SendServerErrorResponse(c, 5101, err)
		return
	}
	if user.Password != reqUser.Password {
		u.SendCustomResponseByBacked(c, "密码错误", "Password error", err)
		return
	}
	token, err := middleware.GenerateToken(user.Name)
	if err != nil {
		u.SendCustomResponseByBacked(c, "生成token失败", "token create error", err)
		return
	}
	u.SendSuccessResponse(c, struct {
		Token string            `json:"token"`
		User  *authcCenter.User `json:"user"`
	}{
		token,
		user,
	})
}
