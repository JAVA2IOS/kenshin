package controllers

import (
	kenshinUtil "kenshin/util"
	"net/http"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

// 请求回调返回参数
type SweetyResponse struct {
	Code    int
	Message string
	Data    interface{}
}

/* 所有的controller都继承当前controller */
type BaseController struct {
	beego.Controller
}

func (baseController *BaseController) Prepare() {
	baseController.Data["SiteName"] = "0NE DATE"
	baseController.Data["AppName"] = "0NEDATE	FILE"
	logs.Info("当前参数:", baseController.GetSession("uid"))
	baseController.Data["Uid"] = baseController.GetSession("uid")

	// println("当前Controller: ", baseController.Controller.Ctx.Input.RunController.Name())
}

func (c *BaseController) Get() {

	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.Data["CustomText"] = "你好，第一个beego web 项目"

	c.TplName = "main.html"
	c.Render()
}

// 成功回调参数
func (c *BaseController) Success(data interface{}) {

	sweetyResponse := SweetyResponse{
		Code:    http.StatusOK,
		Message: "",
		Data:    data,
	}

	c.Data["json"] = sweetyResponse

	c.ServeJSON()
}

// 失败回调
func (c *BaseController) Failure(code int, message string) {

	sweetyResponse := SweetyResponse{
		Code:    code,
		Message: message,
		Data:    nil,
	}

	c.Data["json"] = sweetyResponse

	c.ServeJSON()
}

// 登录认证
func (c *BaseController) Auth() {

	c.TplName = "index.html"
	c.Render()
}

// 登录跳转
func (c *BaseController) ValidateAuth() {

	uid := c.GetString("uid")

	uid = strings.TrimSpace(uid)

	if len(uid) == 0 {
		c.Failure(500, "请输入用户名")
		return
	}

	// 创建根目录

	go func() {
		fileDirectory := "/file/tmp/" + strings.ReplaceAll(uid, " ", "_")

		_, err := kenshinUtil.CreateFileDirectory(fileDirectory)

		if err != nil {
			logs.Warn("目录创建失败:", err.Error())
		} else {
			logs.Info("目录创建成功:", fileDirectory)
		}
	}()

	c.SetSession("uid", uid)

	c.Ctx.Redirect(302, "/")
}
