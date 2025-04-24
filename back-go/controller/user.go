package controller

import (
	"errors"
	"gin-web/middleware"
	"gin-web/models/authcCenter"
	"gin-web/utils"
	"gin-web/utils/extendController"
	"net/http"
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

func (u *UserController) Add(c *gin.Context) {
	var reqUser userRequest
	if err := c.BindJSON(&reqUser); err != nil {
		u.SendParameterErrorResponse(c, err)
		return
	}
	isExist, err := reqUser.User.IsExist()
	if isExist || err != nil {
		u.SendDataDuplicationResponse(c, err)
		return
	}
	if err = reqUser.User.Add(reqUser.AddRoles); err != nil {
		u.SendCustomResponse(c, "添加user失败", "add user failed", err)
		return
	}
	u.SendSuccessResponse(c, "success")
}

func (u *UserController) Deleted(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		u.SendParameterErrorResponse(c, errors.New("参数错误,id为空"))
		return
	}
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		u.SendParameterErrorResponse(c, errors.New("id转化错误为空"))
		return
	}
	if err = new(authcCenter.User).Deleted(idInt64); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	u.SendSuccessResponse(c, "success")
}

func (u *UserController) Update(c *gin.Context) {
	var reqUser userRequest
	if err := c.ShouldBind(&reqUser); err != nil {
		u.SendParameterErrorResponse(c, err)
		return
	}
	if err := reqUser.User.Update(reqUser.AddRoles, reqUser.DeletedRoles); err != nil {
		u.SendCustomResponse(c, "更新user失败", "update user failed", err)
		return
	}
	u.SendSuccessResponse(c, "success")
}

func (u *UserController) GetAll(c *gin.Context) {
	var reqUser authcCenter.User
	reqUser.Name = c.Query("name")
	currPage := c.DefaultQuery("currPage", "1")
	pageSize := c.DefaultQuery("pageSize", "10")
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")
	skip, limit, err := utils.GetPage(currPage, pageSize)
	if err != nil {
		u.SendParameterErrorResponse(c, err)
		return
	}
	resDB, err := reqUser.GetAll(skip, limit, startTime, endTime)
	if err != nil {
		u.SendCustomResponse(c, "查询user失败", "find user failed", err)
		return
	}
	u.SendSuccessResponse(c, resDB)
}

func (u *UserController) Login(c *gin.Context) {
	var reqUser struct {
		Account  string `form:"account"`
		Password string `form:"password"`
	}
	if err := c.ShouldBind(&reqUser); err != nil {
		u.SendParameterErrorResponse(c, err)
		return
	}
	user, err := new(authcCenter.User).GetOne(reqUser.Account, "")
	if err != nil {
		u.SendCustomResponse(c, "该用户不存在", "The user does not exist", err)
		return
	}
	if user.Password != reqUser.Password {
		u.SendCustomResponse(c, "密码错误", "Password error", err)
		return
	}
	token, err := middleware.GenerateToken(user.Name)
	if err != nil {
		u.SendCustomResponse(c, "生成token失败", "token create error", err)
		return
	}
	u.SendSuccessResponse(c, struct {
		Token string `json:"token"`
		User  *authcCenter.User `json:"user"`
	}{
		token,
		user,
	})
}
