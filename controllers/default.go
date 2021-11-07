package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
)

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
}
