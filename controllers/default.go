package controllers

import (
	"net/http"

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

	// println("当前Controller: ", baseController.Controller.Ctx.Input.RunController.Name())
}

func (c *BaseController) Get() {

	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.Data["CustomText"] = "你好，第一个beego web 项目"
	// c.Layout = "main.html"
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
