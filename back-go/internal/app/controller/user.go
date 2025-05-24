package controller

import (
	"fmt"
	"gin-web/internal/app/models/authCenter"
	mysqlDB "gin-web/internal/initialize/mysql"
	"gin-web/internal/middleware"
	"gin-web/pkg"
	"gin-web/pkg/extendController"
	"math/rand"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type userRequest struct {
	User         authCenter.User `json:"user"`
	AddRoles     []int           `json:"addRoles"`
	DeletedRoles []int           `json:"deletedRoles"`
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
	//是否重复
	isExist, err := reqUser.User.IsExist()
	if isExist || err != nil {
		u.SendServerErrorResponse(c, 5101, err)
		return
	}
	rand.Seed(time.Now().Unix()) //根据时间戳生成种子
	//生成盐
	salt := strconv.FormatInt(rand.Int63(), 10) //生成盐
	encryptedPass, err := bcrypt.GenerateFromPassword([]byte(reqUser.User.Password+salt), bcrypt.DefaultCost)
	if err != nil {
		u.SendCustomResponseByBacked(c, "加密失败", "Encryption failed", err)
		return
	}
	reqUser.User.Password, reqUser.User.Salt = string(encryptedPass), salt
	if err = reqUser.User.Insert(reqUser.AddRoles); err != nil {
		u.SendServerErrorResponse(c, 5100, err)
		return
	}
	u.SendSuccessResponse(c, "success")
}

func (u *UserController) Remove(c *gin.Context) {
	//参数校验
	var req struct {
		Id int64 `json:"id"`
	}
	if err := c.ShouldBind(&req); err != nil {
		u.SendParameterErrorResponse(c, 4002, err)
		return
	}
	//DB操作
	if err := new(authCenter.User).Remove(req.Id); err != nil {
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
	fmt.Println("TESTING", reqUser.AddRoles, reqUser.DeletedRoles)
	//DB操作
	if err := reqUser.User.Edit(reqUser.AddRoles, reqUser.DeletedRoles); err != nil {
		u.SendServerErrorResponse(c, 5120, err)
		return
	}
	u.SendSuccessResponse(c, "success")
}

func (u *UserController) GetList(c *gin.Context) {
	//参数校验
	var user authCenter.User
	user.Name, user.Email = c.Query("name"), c.Query("email")
	currPage, pageSize := c.DefaultQuery("currPage", "1"), c.DefaultQuery("pageSize", "10")
	startTime, endTime := c.Query("startTime"), c.Query("endTime")
	skip, limit, err := pkg.GetPage(currPage, pageSize)
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
		Logs  []authCenter.User `json:"users"`
		Total int64             `json:"total"`
	}{
		resDB,
		count,
	})
}

func (u *UserController) GetRolesByUserID(c *gin.Context) {
	userId := c.Query("id")
	if userId == "" {
		u.SendParameterErrorResponse(c, 4001, nil)
		return
	}
	id, err := strconv.Atoi(userId)
	if err != nil {
		u.SendParameterErrorResponse(c, 4003, err)
		return
	}
	var user authCenter.User
	err = mysqlDB.DB.
		Preload("Roles", func(db *gorm.DB) *gorm.DB {
			return db.Select("role.id", "role.name") // 指定表名.字段更稳妥
		}).
		Where("id = ?", id).
		First(&user).Error
	if err != nil {
		u.SendServerErrorResponse(c, 5130, err)
		return
	}
	u.SendSuccessResponse(c, user.Roles)
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
	user := &authCenter.User{
		Account: reqUser.Account,
		Email:   reqUser.Account, // 如果你支持用 email 登录
	}
	user, err := user.GetOne()
	if err != nil {
		u.SendServerErrorResponse(c, 5101, err)
		return
	}
	//校验密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqUser.Password+user.Salt)); err != nil {
		u.SendCustomResponseByBacked(c, "密码错误", "Password error", err)
		return
	}
	token, err := middleware.GenerateToken(user.Name)
	if err != nil {
		u.SendCustomResponseByBacked(c, "生成token失败", "token create error", err)
		return
	}
	u.SendSuccessResponse(c, struct {
		Token string           `json:"token"`
		User  *authCenter.User `json:"user"`
	}{
		token,
		user,
	})
}
