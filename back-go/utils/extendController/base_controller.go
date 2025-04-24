package extendController

import (
	"fmt"
	"gin-web/initialize/runLog"
	"github.com/gin-gonic/gin"
	"net/http"
)

type BaseController struct {
}

func (b *BaseController) SendResponse(c *gin.Context, httpResponseCode, CustomCode int, msg ResponseMsg, data interface{}, err error) {
	if err != nil {
		// 追加错误信息到原有的 msg.ZhCn 和 msg.EnUs 上
		msg.ZhCn = fmt.Sprintf("%s : %v", msg.ZhCn, err)
		msg.EnUs = fmt.Sprintf("%s : %v", msg.EnUs, err)
		runLog.ZapLog.Info(msg.ZhCn + "/|^_^|/" + msg.EnUs) //定义日志输出格式
	}
	c.JSON(httpResponseCode, Response{ //http码
		Code:    CustomCode,
		Message: msg,  //错误信息还是正确信息
		Data:    data, //数据
	})
}

// 成功200
func (b *BaseController) SendSuccessResponse(c *gin.Context, data interface{}) {
	b.SendResponse(c, http.StatusOK, Normal, ResponseMsg{
		ZhCn: "请求成功",
		EnUs: "success",
	}, data, nil)

}

// 自定义错误
func (b *BaseController) SendCustomResponse(c *gin.Context, ZhCn, EnUs string, err error) {
	b.SendResponse(c, http.StatusOK, Normal, ResponseMsg{
		ZhCn: ZhCn,
		EnUs: EnUs,
	}, nil, err)
}

// 参数错误400
func (b *BaseController) SendParameterErrorResponse(c *gin.Context, err error) {
	b.SendResponse(c, http.StatusBadRequest, ParameterError, ResponseMsg{
		ZhCn: "参数错误",
		EnUs: "parameter error",
	}, nil, err)
}

// 权限问题401
func (b *BaseController) SendUnAuthResponse(c *gin.Context) {
	b.SendResponse(c, http.StatusUnauthorized, Unauthorized, ResponseMsg{
		ZhCn: "身份信息不通过",
		EnUs: "Identity information not passed",
	}, nil, nil)
}

// 文件不存在404
func (b *BaseController) SendNotFoundResponse(c *gin.Context, err error) {
	b.SendResponse(c, http.StatusNotFound, NotFound, ResponseMsg{
		ZhCn: "文件不存在",
		EnUs: "File does not exist",
	}, nil, err)
}

// 方法不允许405
func (b *BaseController) SendMethodNotAllowedResponse(c *gin.Context) {
	b.SendResponse(c, http.StatusMethodNotAllowed, NotFound, ResponseMsg{
		ZhCn: "方法不允许",
		EnUs: "Method not allow",
	}, nil, nil)
}

// 重复问题409
func (b *BaseController) SendDataDuplicationResponse(c *gin.Context, err error) {
	b.SendResponse(c, http.StatusConflict, Unauthorized, ResponseMsg{
		ZhCn: "数据重复",
		EnUs: "Data duplication",
	}, nil, err)
}
